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

type LoanHandler struct {
	db *database.Database
}

func NewLoanHandler(db *database.Database) *LoanHandler {
	return &LoanHandler{db: db}
}

func (h *LoanHandler) GetLoans(c *fiber.Ctx) error {
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
		return c.JSON([]models.Loan{})
	}

	// Получаем кредиты для компаний пользователя
	loansCollection := h.db.DB.Collection("loans")
	filter := bson.M{"company_id": bson.M{"$in": companyIDs}}

	// Фильтр по компании если указан
	if companyID := c.Query("company_id"); companyID != "" {
		companyObjectID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Неверный ID компании"})
		}
		filter["company_id"] = companyObjectID
	}

	cursor, err := loansCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения кредитов"})
	}
	defer cursor.Close(context.TODO())

	var loans []models.Loan
	if err = cursor.All(context.TODO(), &loans); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования данных"})
	}

	return c.JSON(loans)
}

func (h *LoanHandler) CreateLoan(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	var loan models.Loan
	if err := c.BodyParser(&loan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// Проверяем что компания принадлежит пользователю
	companiesCollection := h.db.DB.Collection("companies")
	var company models.Company
	err = companiesCollection.FindOne(context.TODO(), bson.M{
		"_id":     loan.CompanyID,
		"user_id": userObjectID,
	}).Decode(&company)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Компания не найдена или нет доступа"})
	}

	// Рассчитываем месячный платеж
	loan.MonthlyPayment = utils.CalculateMonthlyPayment(
		loan.PrincipalAmount,
		loan.InterestRate,
		loan.TermMonths,
	)
	loan.RemainingBalance = loan.PrincipalAmount
	loan.CreatedAt = time.Now()
	loan.UpdatedAt = time.Now()

	loansCollection := h.db.DB.Collection("loans")
	result, err := loansCollection.InsertOne(context.TODO(), loan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка создания кредита"})
	}

	loan.ID = result.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(loan)
}

func (h *LoanHandler) UpdateLoan(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	loanID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID кредита"})
	}

	var loan models.Loan
	if err := c.BodyParser(&loan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// Проверяем что компания принадлежит пользователю
	companiesCollection := h.db.DB.Collection("companies")
	var company models.Company
	err = companiesCollection.FindOne(context.TODO(), bson.M{
		"_id":     loan.CompanyID,
		"user_id": userObjectID,
	}).Decode(&company)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Компания не найдена или нет доступа"})
	}

	// Пересчитываем месячный платеж если изменились параметры
	loan.MonthlyPayment = utils.CalculateMonthlyPayment(
		loan.PrincipalAmount,
		loan.InterestRate,
		loan.TermMonths,
	)
	loan.UpdatedAt = time.Now()

	loansCollection := h.db.DB.Collection("loans")
	filter := bson.M{"_id": loanID, "company_id": loan.CompanyID}
	update := bson.M{"$set": loan}

	result, err := loansCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления кредита"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Кредит не найден"})
	}

	loan.ID = loanID
	return c.JSON(loan)
}

func (h *LoanHandler) DeleteLoan(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	loanID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID кредита"})
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

	loansCollection := h.db.DB.Collection("loans")
	filter := bson.M{
		"_id":        loanID,
		"company_id": bson.M{"$in": companyIDs},
	}

	result, err := loansCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления кредита"})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Кредит не найден"})
	}

	return c.JSON(fiber.Map{"message": "Кредит удален"})
}
