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

type VehicleHandler struct {
	db *database.Database
}

func NewVehicleHandler(db *database.Database) *VehicleHandler {
	return &VehicleHandler{db: db}
}

func (h *VehicleHandler) GetVehicles(c *fiber.Ctx) error {
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

	// Собираем ID компаний
	var companyIDs []primitive.ObjectID
	for _, company := range companies {
		companyIDs = append(companyIDs, company.ID)
	}

	// Если у пользователя нет компаний, возвращаем пустой массив
	if len(companyIDs) == 0 {
		return c.JSON([]models.Vehicle{})
	}

	// Получаем транспорт для компаний пользователя
	vehiclesCollection := h.db.DB.Collection("vehicles")
	filter := bson.M{"company_id": bson.M{"$in": companyIDs}}

	// Фильтр по компании если указан
	if companyID := c.Query("company_id"); companyID != "" {
		companyObjectID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Неверный ID компании"})
		}
		filter["company_id"] = companyObjectID
	}

	cursor, err := vehiclesCollection.Find(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения транспорта"})
	}
	defer cursor.Close(context.TODO())

	var vehicles []models.Vehicle
	if err = cursor.All(context.TODO(), &vehicles); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка декодирования данных"})
	}

	return c.JSON(vehicles)
}

func (h *VehicleHandler) CreateVehicle(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	var vehicle models.Vehicle
	if err := c.BodyParser(&vehicle); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// Проверяем что компания принадлежит пользователю
	companiesCollection := h.db.DB.Collection("companies")
	var company models.Company
	err = companiesCollection.FindOne(context.TODO(), bson.M{
		"_id":     vehicle.CompanyID,
		"user_id": userObjectID,
	}).Decode(&company)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Компания не найдена или нет доступа"})
	}

	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()

	vehiclesCollection := h.db.DB.Collection("vehicles")
	result, err := vehiclesCollection.InsertOne(context.TODO(), vehicle)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка создания транспорта"})
	}

	vehicle.ID = result.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(vehicle)
}

func (h *VehicleHandler) UpdateVehicle(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	vehicleID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID транспорта"})
	}

	var vehicle models.Vehicle
	if err := c.BodyParser(&vehicle); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// Проверяем что компания принадлежит пользователю
	companiesCollection := h.db.DB.Collection("companies")
	var company models.Company
	err = companiesCollection.FindOne(context.TODO(), bson.M{
		"_id":     vehicle.CompanyID,
		"user_id": userObjectID,
	}).Decode(&company)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Компания не найдена или нет доступа"})
	}

	vehicle.UpdatedAt = time.Now()

	vehiclesCollection := h.db.DB.Collection("vehicles")
	filter := bson.M{"_id": vehicleID, "company_id": vehicle.CompanyID}
	update := bson.M{"$set": vehicle}

	result, err := vehiclesCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления транспорта"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Транспорт не найден"})
	}

	vehicle.ID = vehicleID
	return c.JSON(vehicle)
}

func (h *VehicleHandler) DeleteVehicle(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	vehicleID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID транспорта"})
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

	vehiclesCollection := h.db.DB.Collection("vehicles")
	filter := bson.M{
		"_id":        vehicleID,
		"company_id": bson.M{"$in": companyIDs},
	}

	result, err := vehiclesCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления транспорта"})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Транспорт не найден"})
	}

	return c.JSON(fiber.Map{"message": "Транспорт удален"})
}
