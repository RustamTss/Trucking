package handlers

import (
	"business-schedule-backend/database"
	"business-schedule-backend/models"
	"business-schedule-backend/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	db        *database.Database
	jwtSecret string
}

func NewAuthHandler(db *database.Database, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// Проверяем существует ли пользователь
	collection := h.db.DB.Collection("users")
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Пользователь уже существует"})
	}

	// Хешируем пароль
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка хеширования пароля"})
	}

	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Сохраняем пользователя
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка создания пользователя"})
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	user.Password = "" // Не возвращаем пароль

	// Генерируем JWT токен
	token, err := utils.GenerateJWT(user.ID.Hex(), h.jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка генерации токена"})
	}

	return c.JSON(models.LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginReq models.LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// Ищем пользователя
	collection := h.db.DB.Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": loginReq.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(401).JSON(fiber.Map{"error": "Неверные учетные данные"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка поиска пользователя"})
	}

	// Проверяем пароль
	if !utils.CheckPasswordHash(loginReq.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Неверные учетные данные"})
	}

	// Генерируем JWT токен
	token, err := utils.GenerateJWT(user.ID.Hex(), h.jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка генерации токена"})
	}

	user.Password = "" // Не возвращаем пароль

	return c.JSON(models.LoginResponse{
		Token: token,
		User:  user,
	})
}
