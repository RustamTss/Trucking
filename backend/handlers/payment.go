package handlers

import (
	"business-schedule-backend/database"
	"business-schedule-backend/middleware"
	"business-schedule-backend/models"
	"business-schedule-backend/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentHandler struct {
	db *database.Database
}

func NewPaymentHandler(db *database.Database) *PaymentHandler {
	return &PaymentHandler{db: db}
}

func (h *PaymentHandler) GetPaymentsByLoan(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	loanID, err := primitive.ObjectIDFromHex(c.Params("loanId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID кредита"})
	}

	// Проверяем что кредит принадлежит пользователю
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

	loansCollection := h.db.DB.Collection("loans")
	var loan models.Loan
	err = loansCollection.FindOne(context.TODO(), bson.M{
		"_id":        loanID,
		"company_id": bson.M{"$in": companyIDs},
	}).Decode(&loan)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Кредит не найден"})
	}

	// Получаем платежи по кредиту
	paymentsCollection := h.db.DB.Collection("payments")
	cursor, err := paymentsCollection.Find(context.TODO(), bson.M{"loan_id": loanID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения платежей"})
	}
	defer cursor.Close(context.TODO())

	var payments []models.Payment
	if err = cursor.All(context.TODO(), &payments); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования данных"})
	}

	return c.JSON(payments)
}

func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	var payment models.Payment
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// Проверяем что кредит принадлежит пользователю
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

	loansCollection := h.db.DB.Collection("loans")
	var loan models.Loan
	err = loansCollection.FindOne(context.TODO(), bson.M{
		"_id":        payment.LoanID,
		"company_id": bson.M{"$in": companyIDs},
	}).Decode(&loan)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Кредит не найден"})
	}

	// Рассчитываем процентную часть платежа
	interestPayment := utils.CalculateInterestPayment(loan.RemainingBalance, loan.InterestRate)
	principalPayment := payment.TotalPaid - interestPayment

	if principalPayment < 0 {
		principalPayment = 0
		interestPayment = payment.TotalPaid
	}

	payment.PrincipalPaid = principalPayment
	payment.InterestPaid = interestPayment
	payment.RemainingBalance = utils.CalculateRemainingBalance(loan.RemainingBalance, principalPayment)
	payment.CreatedAt = time.Now()

	// Сохраняем платеж
	paymentsCollection := h.db.DB.Collection("payments")
	result, err := paymentsCollection.InsertOne(context.TODO(), payment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка создания платежа"})
	}

	payment.ID = result.InsertedID.(primitive.ObjectID)

	// Обновляем остаток по кредиту
	loan.RemainingBalance = payment.RemainingBalance
	if loan.RemainingBalance <= 0 {
		loan.Status = "paid_off"
	}
	loan.UpdatedAt = time.Now()

	_, err = loansCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": loan.ID},
		bson.M{"$set": bson.M{
			"remaining_balance": loan.RemainingBalance,
			"status":            loan.Status,
			"updated_at":        loan.UpdatedAt,
		}},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления кредита"})
	}

	return c.Status(201).JSON(payment)
}

func (h *PaymentHandler) GetPayments(c *fiber.Ctx) error {
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

	// Если у пользователя нет компаний, возвращаем пустой массив
	if len(companyIDs) == 0 {
		return c.JSON([]models.Payment{})
	}

	// Получаем все кредиты пользователя
	loansCollection := h.db.DB.Collection("loans")
	loansCursor, err := loansCollection.Find(context.TODO(), bson.M{"company_id": bson.M{"$in": companyIDs}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения кредитов"})
	}
	defer loansCursor.Close(context.TODO())

	var loans []models.Loan
	if err = loansCursor.All(context.TODO(), &loans); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования кредитов"})
	}

	var loanIDs []primitive.ObjectID
	for _, loan := range loans {
		loanIDs = append(loanIDs, loan.ID)
	}

	// Получаем все платежи по кредитам пользователя
	paymentsCollection := h.db.DB.Collection("payments")
	cursor, err := paymentsCollection.Find(context.TODO(), bson.M{"loan_id": bson.M{"$in": loanIDs}})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения платежей"})
	}
	defer cursor.Close(context.TODO())

	var payments []models.Payment
	if err = cursor.All(context.TODO(), &payments); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования данных"})
	}

	return c.JSON(payments)
}
