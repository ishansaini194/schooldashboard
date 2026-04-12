package app

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/handlers"
	"github.com/ishansaini194/dashboard/internal/server"
	"github.com/ishansaini194/dashboard/middleware"
)

func New() *server.Server {
	database.Connect("/app/data/school.db")
	database.Run(database.DB)

	srv := server.New()

	srv.App.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://127.0.0.1:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	api := srv.App.Group("/api")

	// ── Public routes (no auth) ──
	api.Post("/auth/register", handlers.Register)
	api.Post("/auth/login", handlers.Login)

	// ── Protected routes (auth required) ──
	protected := api.Group("/", middleware.AuthRequired())

	// Student routes
	protected.Get("/students/class/:class", handlers.GetStudents)
	protected.Get("/students/:roll_no", handlers.GetStudent)
	protected.Post("/students", handlers.CreateStudent)
	protected.Put("/students/:roll_no", handlers.UpdateStudent)
	protected.Delete("/students/:roll_no", handlers.DeleteStudent)

	// Class routes
	protected.Get("/classes", handlers.GetClasses)
	protected.Get("/classes/:id", handlers.GetClass)
	protected.Post("/classes", handlers.CreateClass)
	protected.Put("/classes/:id", handlers.UpdateClass)
	protected.Delete("/classes/:id", handlers.DeleteClass)

	// Fee routes
	protected.Post("/fees/pay", handlers.PayFee)
	protected.Get("/fees/student/:student_id", handlers.GetStudentFees)
	protected.Get("/fees/student/:student_id/yearly", handlers.GetStudentYearlySummary)
	protected.Get("/fees/class/:class/month/:month/year/:year", handlers.GetClassFeeStatus)
	protected.Get("/fees/pending/:class/:month/:year", handlers.GetPendingFees)
	protected.Get("/fees/pending/all", handlers.GetAllPendingFees)
	protected.Get("/fees/receipt/:receipt_no", handlers.GetReceipt)
	protected.Get("/fees/recent", handlers.GetRecentPayments)
	protected.Get("/fees/overdue", handlers.GetOverdueFees)
	protected.Put("/fees/:id/complete", handlers.CompleteFee)

	// Dashboard routes
	protected.Get("/dashboard/summary", handlers.GetDashboardSummary)

	return srv
}
