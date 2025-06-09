package handlers

import (
	"business-schedule-backend/database"
	"business-schedule-backend/middleware"
	"business-schedule-backend/models"
	"business-schedule-backend/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduleHandler struct {
	db *database.Database
}

func NewScheduleHandler(db *database.Database) *ScheduleHandler {
	return &ScheduleHandler{db: db}
}

type DebtScheduleItem struct {
	CompanyName    string  `json:"company_name"`
	TotalDebt      float64 `json:"total_debt"`
	MonthlyPayment float64 `json:"monthly_payment"`
	VehiclesCount  int     `json:"vehicles_count"`
}

type AmortizationScheduleItem struct {
	PaymentNumber    int     `json:"payment_number"`
	PaymentDate      string  `json:"payment_date"`
	PrincipalPayment float64 `json:"principal_payment"`
	InterestPayment  float64 `json:"interest_payment"`
	TotalPayment     float64 `json:"total_payment"`
	RemainingBalance float64 `json:"remaining_balance"`
}

type DepreciationScheduleItem struct {
	VehicleID          string  `json:"vehicle_id"`
	VehicleName        string  `json:"vehicle_name"`
	PurchasePrice      float64 `json:"purchase_price"`
	CurrentValue       float64 `json:"current_value"`
	DepreciationAmount float64 `json:"depreciation_amount"`
	AgeYears           float64 `json:"age_years"`
}

func (h *ScheduleHandler) GetDebtSchedule(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	// Получаем компании пользователя
	companiesCollection := h.db.DB.Collection("companies")
	companiesCursor, err := companiesCollection.Find(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения компаний"})
	}
	defer companiesCursor.Close(context.TODO())

	var companies []models.Company
	if err = companiesCursor.All(context.TODO(), &companies); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования компаний"})
	}

	var debtSchedule []DebtScheduleItem

	for _, company := range companies {
		// Получаем кредиты компании
		loansCollection := h.db.DB.Collection("loans")
		loansCursor, err := loansCollection.Find(context.TODO(), bson.M{
			"company_id": company.ID,
			"status":     "active",
		})
		if err != nil {
			continue
		}

		var loans []models.Loan
		if err = loansCursor.All(context.TODO(), &loans); err != nil {
			loansCursor.Close(context.TODO())
			continue
		}
		loansCursor.Close(context.TODO())

		// Получаем транспорт компании
		vehiclesCollection := h.db.DB.Collection("vehicles")
		vehiclesCount, err := vehiclesCollection.CountDocuments(context.TODO(), bson.M{
			"company_id": company.ID,
			"status":     "active",
		})
		if err != nil {
			vehiclesCount = 0
		}

		// Рассчитываем общий долг и платежи
		totalDebt := 0.0
		monthlyPayment := 0.0
		for _, loan := range loans {
			totalDebt += loan.RemainingBalance
			monthlyPayment += loan.MonthlyPayment
		}

		debtSchedule = append(debtSchedule, DebtScheduleItem{
			CompanyName:    company.Name,
			TotalDebt:      totalDebt,
			MonthlyPayment: monthlyPayment,
			VehiclesCount:  int(vehiclesCount),
		})
	}

	return c.JSON(debtSchedule)
}

func (h *ScheduleHandler) GetAmortizationSchedule(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	// Получаем компании пользователя
	companiesCollection := h.db.DB.Collection("companies")
	companiesCursor, err := companiesCollection.Find(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения компаний"})
	}
	defer companiesCursor.Close(context.TODO())

	var companies []models.Company
	if err = companiesCursor.All(context.TODO(), &companies); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования компаний"})
	}

	var companyIDs []primitive.ObjectID
	for _, company := range companies {
		companyIDs = append(companyIDs, company.ID)
	}

	// Фильтр по кредиту если указан
	filter := bson.M{"company_id": bson.M{"$in": companyIDs}, "status": "active"}
	if loanID := c.Query("loan_id"); loanID != "" {
		loanObjectID, err := primitive.ObjectIDFromHex(loanID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Неверный ID кредита"})
		}
		filter["_id"] = loanObjectID
	}

	// Получаем кредиты
	loansCollection := h.db.DB.Collection("loans")
	loansCursor, err := loansCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения кредитов"})
	}
	defer loansCursor.Close(context.TODO())

	var loans []models.Loan
	if err = loansCursor.All(context.TODO(), &loans); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования кредитов"})
	}

	var amortizationSchedule []AmortizationScheduleItem

	for _, loan := range loans {
		balance := loan.RemainingBalance
		monthlyRate := loan.InterestRate / 100 / 12

		for i := 1; i <= loan.TermMonths && balance > 0; i++ {
			interestPayment := balance * monthlyRate
			principalPayment := loan.MonthlyPayment - interestPayment

			if principalPayment > balance {
				principalPayment = balance
			}

			balance -= principalPayment

			// Рассчитываем дату платежа
			paymentDate := loan.StartDate.AddDate(0, i-1, 0)

			amortizationSchedule = append(amortizationSchedule, AmortizationScheduleItem{
				PaymentNumber:    i,
				PaymentDate:      paymentDate.Format("2006-01-02"),
				PrincipalPayment: principalPayment,
				InterestPayment:  interestPayment,
				TotalPayment:     principalPayment + interestPayment,
				RemainingBalance: balance,
			})
		}
	}

	return c.JSON(amortizationSchedule)
}

func (h *ScheduleHandler) GetDepreciationSchedule(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	// Получаем компании пользователя
	companiesCollection := h.db.DB.Collection("companies")
	companiesCursor, err := companiesCollection.Find(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения компаний"})
	}
	defer companiesCursor.Close(context.TODO())

	var companies []models.Company
	if err = companiesCursor.All(context.TODO(), &companies); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования компаний"})
	}

	var companyIDs []primitive.ObjectID
	for _, company := range companies {
		companyIDs = append(companyIDs, company.ID)
	}

	// Фильтр по компании если указан
	filter := bson.M{"company_id": bson.M{"$in": companyIDs}}
	if companyID := c.Query("company_id"); companyID != "" {
		companyObjectID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Неверный ID компании"})
		}
		filter["company_id"] = companyObjectID
	}

	// Получаем транспорт
	vehiclesCollection := h.db.DB.Collection("vehicles")
	vehiclesCursor, err := vehiclesCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения транспорта"})
	}
	defer vehiclesCursor.Close(context.TODO())

	var vehicles []models.Vehicle
	if err = vehiclesCursor.All(context.TODO(), &vehicles); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования транспорта"})
	}

	var depreciationSchedule []DepreciationScheduleItem

	for _, vehicle := range vehicles {
		// Рассчитываем возраст транспорта
		ageYears := utils.CalculateVehicleAge(vehicle.PurchaseDate)

		// Используем стандартный срок службы: 10 лет для траков, 15 лет для трейлеров
		usefulLife := 10
		if vehicle.Type == "trailer" {
			usefulLife = 15
		}

		// Рассчитываем амортизацию
		depreciationAmount := utils.CalculateDepreciation(vehicle.PurchasePrice, usefulLife, ageYears)
		currentValue := vehicle.PurchasePrice - depreciationAmount

		if currentValue < 0 {
			currentValue = 0
		}

		vehicleName := vehicle.Make + " " + vehicle.Model + " (" + string(rune(vehicle.Year)) + ")"

		depreciationSchedule = append(depreciationSchedule, DepreciationScheduleItem{
			VehicleID:          vehicle.ID.Hex(),
			VehicleName:        vehicleName,
			PurchasePrice:      vehicle.PurchasePrice,
			CurrentValue:       currentValue,
			DepreciationAmount: depreciationAmount,
			AgeYears:           ageYears,
		})
	}

	return c.JSON(depreciationSchedule)
}

type DashboardStats struct {
	TotalCompanies    int     `json:"total_companies"`
	TotalVehicles     int     `json:"total_vehicles"`
	TotalActiveLoans  int     `json:"total_active_loans"`
	TotalDebt         float64 `json:"total_debt"`
	MonthlyPayments   float64 `json:"monthly_payments"`
	TotalAssetValue   float64 `json:"total_asset_value"`
	TotalPaymentsYear float64 `json:"total_payments_year"`
}

func (h *ScheduleHandler) GetDashboardStats(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	var stats DashboardStats

	// Получаем компании пользователя
	companiesCollection := h.db.DB.Collection("companies")
	companiesCount, err := companiesCollection.CountDocuments(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения компаний"})
	}
	stats.TotalCompanies = int(companiesCount)

	companiesCursor, err := companiesCollection.Find(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения компаний"})
	}
	defer companiesCursor.Close(context.TODO())

	var companies []models.Company
	if err = companiesCursor.All(context.TODO(), &companies); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования компаний"})
	}

	var companyIDs []primitive.ObjectID
	for _, company := range companies {
		companyIDs = append(companyIDs, company.ID)
	}

	if len(companyIDs) == 0 {
		return c.JSON(stats)
	}

	// Получаем транспорт
	vehiclesCollection := h.db.DB.Collection("vehicles")
	vehiclesCount, err := vehiclesCollection.CountDocuments(context.TODO(), bson.M{"company_id": bson.M{"$in": companyIDs}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения транспорта"})
	}
	stats.TotalVehicles = int(vehiclesCount)

	// Получаем активные кредиты и рассчитываем общий долг
	loansCollection := h.db.DB.Collection("loans")
	activeLoansCount, err := loansCollection.CountDocuments(context.TODO(), bson.M{
		"company_id": bson.M{"$in": companyIDs},
		"status":     "active",
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения кредитов"})
	}
	stats.TotalActiveLoans = int(activeLoansCount)

	loansCursor, err := loansCollection.Find(context.TODO(), bson.M{
		"company_id": bson.M{"$in": companyIDs},
		"status":     "active",
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения кредитов"})
	}
	defer loansCursor.Close(context.TODO())

	var loans []models.Loan
	if err = loansCursor.All(context.TODO(), &loans); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования кредитов"})
	}

	for _, loan := range loans {
		stats.TotalDebt += loan.RemainingBalance
		stats.MonthlyPayments += loan.MonthlyPayment
	}

	// Получаем общую стоимость активов (транспорт)
	vehiclesCursor, err := vehiclesCollection.Find(context.TODO(), bson.M{"company_id": bson.M{"$in": companyIDs}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения транспорта"})
	}
	defer vehiclesCursor.Close(context.TODO())

	var vehicles []models.Vehicle
	if err = vehiclesCursor.All(context.TODO(), &vehicles); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования транспорта"})
	}

	for _, vehicle := range vehicles {
		// Рассчитываем текущую стоимость с учетом амортизации
		ageYears := utils.CalculateVehicleAge(vehicle.PurchaseDate)
		usefulLife := 10
		if vehicle.Type == "trailer" {
			usefulLife = 15
		}

		depreciationAmount := utils.CalculateDepreciation(vehicle.PurchasePrice, usefulLife, ageYears)
		currentValue := vehicle.PurchasePrice - depreciationAmount

		if currentValue < 0 {
			currentValue = 0
		}

		stats.TotalAssetValue += currentValue
	}

	// Рассчитываем общие платежи за год
	stats.TotalPaymentsYear = stats.MonthlyPayments * 12

	return c.JSON(stats)
}
