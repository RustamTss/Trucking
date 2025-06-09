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
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	db *database.Database
}

func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{db: db}
}

// GetUsers получает список всех пользователей (только для админов или текущего пользователя)
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	// Получаем только текущего пользователя (можно расширить для админов)
	collection := h.db.DB.Collection("users")
	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
	}

	// Убираем пароль из ответа
	user.Password = ""

	return c.JSON([]models.User{user})
}

// GetUser получает конкретного пользователя по ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	currentUserID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	requestedUserID := c.Params("id")

	// Пользователь может получить только свои данные
	if currentUserID != requestedUserID {
		return c.Status(403).JSON(fiber.Map{"error": "Доступ запрещен"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(requestedUserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	collection := h.db.DB.Collection("users")
	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения пользователя"})
	}

	// Убираем пароль из ответа
	user.Password = ""

	return c.JSON(user)
}

// UpdateUser обновляет данные пользователя
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	currentUserID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	requestedUserID := c.Params("id")

	// Пользователь может обновить только свои данные
	if currentUserID != requestedUserID {
		return c.Status(403).JSON(fiber.Map{"error": "Доступ запрещен"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(requestedUserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	collection := h.db.DB.Collection("users")

	// Проверяем существует ли пользователь
	var existingUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": userObjectID}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения пользователя"})
	}

	// Подготавливаем данные для обновления
	updateFields := bson.M{
		"updated_at": time.Now(),
	}

	// Обновляем только переданные поля
	if updateData.Name != "" {
		updateFields["name"] = updateData.Name
	}

	if updateData.Email != "" {
		// Проверяем уникальность email (если изменился)
		if updateData.Email != existingUser.Email {
			var emailCheck models.User
			err = collection.FindOne(context.TODO(), bson.M{"email": updateData.Email}).Decode(&emailCheck)
			if err == nil {
				return c.Status(400).JSON(fiber.Map{"error": "Email уже используется"})
			}
		}
		updateFields["email"] = updateData.Email
	}

	// Если передан пароль, хешируем его
	if updateData.Password != "" {
		hashedPassword, err := utils.HashPassword(updateData.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка хеширования пароля"})
		}
		updateFields["password"] = hashedPassword
	}

	// Обновляем пользователя
	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userObjectID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления пользователя"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
	}

	// Получаем обновленного пользователя
	var updatedUser models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": userObjectID}).Decode(&updatedUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения обновленного пользователя"})
	}

	// Убираем пароль из ответа
	updatedUser.Password = ""

	return c.JSON(updatedUser)
}

// DeleteUser удаляет пользователя и все его данные
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	currentUserID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	requestedUserID := c.Params("id")

	// Пользователь может удалить только себя
	if currentUserID != requestedUserID {
		return c.Status(403).JSON(fiber.Map{"error": "Доступ запрещен"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(requestedUserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	// Начинаем транзакцию для удаления всех связанных данных
	usersCollection := h.db.DB.Collection("users")
	companiesCollection := h.db.DB.Collection("companies")
	vehiclesCollection := h.db.DB.Collection("vehicles")
	loansCollection := h.db.DB.Collection("loans")
	paymentsCollection := h.db.DB.Collection("payments")

	// Проверяем существует ли пользователь
	var user models.User
	err = usersCollection.FindOne(context.TODO(), bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения пользователя"})
	}

	// Получаем компании пользователя
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

	// Получаем кредиты для удаления платежей
	if len(companyIDs) > 0 {
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

		// Удаляем платежи
		if len(loanIDs) > 0 {
			_, err = paymentsCollection.DeleteMany(context.TODO(), bson.M{"loan_id": bson.M{"$in": loanIDs}})
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления платежей"})
			}
		}

		// Удаляем кредиты
		_, err = loansCollection.DeleteMany(context.TODO(), bson.M{"company_id": bson.M{"$in": companyIDs}})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления кредитов"})
		}

		// Удаляем транспорт
		_, err = vehiclesCollection.DeleteMany(context.TODO(), bson.M{"company_id": bson.M{"$in": companyIDs}})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления транспорта"})
		}
	}

	// Удаляем компании
	_, err = companiesCollection.DeleteMany(context.TODO(), bson.M{"user_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления компаний"})
	}

	// Удаляем пользователя
	result, err := usersCollection.DeleteOne(context.TODO(), bson.M{"_id": userObjectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления пользователя"})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
	}

	return c.JSON(fiber.Map{
		"message":      "Пользователь и все связанные данные успешно удалены",
		"deleted_user": user.Email,
	})
}

// GetProfile получает профиль текущего пользователя
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Не удалось получить ID пользователя"})
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID пользователя"})
	}

	collection := h.db.DB.Collection("users")
	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения пользователя"})
	}

	// Убираем пароль из ответа
	user.Password = ""

	return c.JSON(user)
}
