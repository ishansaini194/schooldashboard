package app

import (
	"os"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/handlers"
	"github.com/ishansaini194/dashboard/internal/server"
	"github.com/ishansaini194/dashboard/middleware"
)

func New() *server.Server {
	database.Connect()
	database.Run(database.DB)

	srv := server.New()

	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000,http://127.0.0.1:3000,http://localhost:5500,http://127.0.0.1:5500"
	}

	srv.App.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	api := srv.App.Group("/api")

	// ── Public ──
	api.Post("/auth/register", handlers.Register)
	api.Post("/auth/login", handlers.Login)

	// ── Protected ──
	protected := api.Group("/", middleware.AuthRequired())

	// Auth
	protected.Put("/auth/change-password", handlers.ChangePassword)
	protected.Put("/auth/reset-password/:user_id", handlers.ResetPassword)
	protected.Get("/users/epunjab/:epunjab_id", handlers.GetUserByEpunjabID)

	// Teachers
	protected.Get("/teachers", handlers.GetTeachers)
	protected.Post("/teachers", handlers.CreateTeacher)
	protected.Delete("/teachers/:id", handlers.DeleteTeacher)

	// Students
	protected.Post("/students", handlers.CreateStudent)
	protected.Get("/students/epunjab/:epunjab_id", handlers.GetStudentByEpunjabID)
	protected.Get("/students/class/:class", handlers.GetStudents)
	protected.Get("/students/:roll_no", handlers.GetStudent)
	protected.Put("/students/:roll_no", handlers.UpdateStudent)
	protected.Delete("/students/:roll_no", handlers.DeleteStudent)

	// Classes
	protected.Get("/classes", handlers.GetClasses)
	protected.Get("/classes/:id", handlers.GetClass)
	protected.Post("/classes", handlers.CreateClass)
	protected.Put("/classes/:id", handlers.UpdateClass)
	protected.Delete("/classes/:id", handlers.DeleteClass)

	// Fees
	protected.Post("/fees/pay", handlers.PayFee)
	protected.Put("/fees/:id/complete", handlers.CompleteFee)
	protected.Get("/fees/class/:class/month/:month/year/:year", handlers.GetClassFeeStatus)
	protected.Get("/fees/pending/all", handlers.GetAllPendingFees)
	protected.Get("/fees/pending/:class/:month/:year", handlers.GetPendingFees)
	protected.Get("/fees/recent", handlers.GetRecentPayments)
	protected.Get("/fees/overdue", handlers.GetOverdueFees)
	protected.Get("/fees/student/:student_id", handlers.GetStudentFees)
	protected.Get("/fees/student/:student_id/yearly", handlers.GetStudentYearlySummary)
	protected.Get("/fees/receipt/:receipt_no", handlers.GetReceipt)

	// Homework
	protected.Post("/homework", handlers.CreateHomework)
	protected.Put("/homework/:id", handlers.UpdateHomework)
	protected.Delete("/homework/:id", handlers.DeleteHomework)
	protected.Get("/homework/class/:class/section/:section", handlers.GetHomework)
	protected.Get("/homework/:id", handlers.GetHomeworkByID)

	// Notices
	protected.Post("/notices", handlers.CreateNotice)
	protected.Put("/notices/:id", handlers.UpdateNotice)
	protected.Delete("/notices/:id", handlers.DeleteNotice)
	protected.Get("/notices", handlers.GetNotices)
	protected.Get("/notices/:id", handlers.GetNoticeByID)

	// Papers
	protected.Post("/papers", handlers.CreatePaper)
	protected.Put("/papers/:id", handlers.UpdatePaper)
	protected.Delete("/papers/:id", handlers.DeletePaper)
	protected.Get("/papers", handlers.GetPapers)
	protected.Get("/papers/:id", handlers.GetPaperByID)

	// Results
	protected.Post("/results", handlers.CreateResult)
	protected.Put("/results/:id", handlers.UpdateResult)
	protected.Delete("/results/:id", handlers.DeleteResult)
	protected.Get("/results/mine", handlers.GetMyResults)
	protected.Get("/results/student/:student_id", handlers.GetStudentResults)
	protected.Get("/results/class/:class/section/:section", handlers.GetClassResults)

	// Dashboard
	protected.Get("/dashboard/summary", handlers.GetDashboardSummary)

	return srv
}
