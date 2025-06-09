package routes

import (
	"business-schedule-backend/database"
	"business-schedule-backend/handlers"
	"business-schedule-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *database.Database, jwtSecret string) {
	// Здоровье приложения
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Business Schedule API v1.0",
		})
	})

	api := app.Group("/api")

	// Аутентификация
	auth := api.Group("/auth")
	authHandler := handlers.NewAuthHandler(db, jwtSecret)
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)

	// Защищенные маршруты
	protected := api.Group("", middleware.JWTMiddleware(jwtSecret))

	// Компании
	companies := protected.Group("/companies")
	companyHandler := handlers.NewCompanyHandler(db)
	companies.Get("/", companyHandler.GetCompanies)
	companies.Post("/", companyHandler.CreateCompany)
	companies.Put("/:id", companyHandler.UpdateCompany)
	companies.Delete("/:id", companyHandler.DeleteCompany)

	// Транспорт
	vehicles := protected.Group("/vehicles")
	vehicleHandler := handlers.NewVehicleHandler(db)
	vehicles.Get("/", vehicleHandler.GetVehicles)
	vehicles.Post("/", vehicleHandler.CreateVehicle)
	vehicles.Put("/:id", vehicleHandler.UpdateVehicle)
	vehicles.Delete("/:id", vehicleHandler.DeleteVehicle)

	// Кредиты
	loans := protected.Group("/loans")
	loanHandler := handlers.NewLoanHandler(db)
	loans.Get("/", loanHandler.GetLoans)
	loans.Post("/", loanHandler.CreateLoan)
	loans.Put("/:id", loanHandler.UpdateLoan)
	loans.Delete("/:id", loanHandler.DeleteLoan)

	// Платежи
	payments := protected.Group("/payments")
	paymentHandler := handlers.NewPaymentHandler(db)
	payments.Get("/", paymentHandler.GetPayments) // Все платежи пользователя
	payments.Get("/loan/:loanId", paymentHandler.GetPaymentsByLoan)
	payments.Post("/", paymentHandler.CreatePayment)

	// Финансовые отчеты
	schedules := protected.Group("/schedules")
	scheduleHandler := handlers.NewScheduleHandler(db)
	schedules.Get("/debt", scheduleHandler.GetDebtSchedule)
	schedules.Get("/amortization", scheduleHandler.GetAmortizationSchedule)
	schedules.Get("/depreciation", scheduleHandler.GetDepreciationSchedule)

	// Статистика
	stats := protected.Group("/stats")
	stats.Get("/dashboard", scheduleHandler.GetDashboardStats)

	// Пользователи
	users := protected.Group("/users")
	userHandler := handlers.NewUserHandler(db)
	users.Get("/", userHandler.GetUsers)
	users.Get("/profile", userHandler.GetProfile)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
