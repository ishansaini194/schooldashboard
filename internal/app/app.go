package app

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/handlers"
	"github.com/ishansaini194/dashboard/internal/server"
)

func New() *server.Server {
	database.Connect("/app/data/school.db")
	database.Run(database.DB)

	srv := server.New()

	srv.App.Use(cors.New((cors.Config{
		AllowOrigins: "http://localhost:3000, http://127.0.0.1:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	})))

	api := srv.App.Group("/api")

	// Student routes
	api.Get("/students/class/:class", handlers.GetStudents)
	api.Get("/students/:roll_no", handlers.GetStudent)
	api.Post("/students", handlers.CreateStudent)
	api.Put("/students/:roll_no", handlers.UpdateStudent)
	api.Delete("/students/:roll_no", handlers.DeleteStudent)

	// Class routes
	api.Get("/classes", handlers.GetClasses)
	api.Get("/classes/:id", handlers.GetClass)
	api.Post("/classes", handlers.CreateClass)
	api.Put("/classes/:id", handlers.UpdateClass)
	api.Delete("/classes/:id", handlers.DeleteClass)

	// fee routes
	api.Post("/fees/pay", handlers.PayFee)
	api.Get("/fees/student/:student_id", handlers.GetStudentFees)
	api.Get("/fees/class/:class/month/:month/year/:year", handlers.GetClassFeeStatus)
	api.Get("/fees/pending/:class/:month/:year", handlers.GetPendingFees)
	api.Get("/fees/receipt/:receipt_no", handlers.GetReceipt)
	api.Get("/fees/student/:student_id/yearly", handlers.GetStudentYearlySummary)
	api.Put("/fees/:id/complete", handlers.CompleteFee)

	// dashboard routes
	api.Get("/dashboard/summary", handlers.GetDashboardSummary)
	api.Get("/fees/recent", handlers.GetRecentPayments)
	api.Get("/fees/overdue", handlers.GetOverdueFees)
	api.Get("/fees/pending/all", handlers.GetAllPendingFees)

	return srv
}
