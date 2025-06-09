package handlers

import (
	"business-schedule-backend/database"
	"business-schedule-backend/middleware"
	"business-schedule-backend/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyHandler struct {
	db *database.Database
}

func NewCompanyHandler(db *database.Database) *CompanyHandler {
	return &CompanyHandler{db: db}
}

func (h *CompanyHandler) GetCompanies(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	collection := h.db.DB.Collection("companies")
	cursor, err := collection.Find(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения компаний"})
	}
	defer cursor.Close(context.TODO())

	var companies []models.Company
	if err = cursor.All(context.TODO(), &companies); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования данных"})
	}

	return c.JSON(companies)
}

func (h *CompanyHandler) CreateCompany(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	var company models.Company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	company.UserID = userObjectID
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	collection := h.db.DB.Collection("companies")
	result, err := collection.InsertOne(context.TODO(), company)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка создания компании"})
	}

	company.ID = result.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(company)
}

func (h *CompanyHandler) UpdateCompany(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	companyID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID компании"})
	}

	var company models.Company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	company.UpdatedAt = time.Now()

	collection := h.db.DB.Collection("companies")
	filter := bson.M{"_id": companyID, "user_id": userObjectID}
	update := bson.M{"$set": company}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления компании"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Компания не найдена"})
	}

	company.ID = companyID
	return c.JSON(company)
}

func (h *CompanyHandler) DeleteCompany(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	companyID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID компании"})
	}

	collection := h.db.DB.Collection("companies")
	filter := bson.M{"_id": companyID, "user_id": userObjectID}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления компании"})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Компания не найдена"})
	}

	return c.JSON(fiber.Map{"message": "Компания удалена"})
}
