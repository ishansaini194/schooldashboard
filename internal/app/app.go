package app

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ishansaini194/dashboard/database"
	"github.com/ishansaini194/dashboard/handlers"
	"github.com/ishansaini194/dashboard/internal/server"
)

func New() *server.Server {
	database.Connect("school.db")
	database.Run(database.DB)

	srv := server.New()

	srv.App.Use(cors.New())

	api := srv.App.Group("/api")

	api.Get("/students/class/:class_id", handlers.GetStudents)
	api.Get("/students/:roll_no", handlers.GetStudent)
	api.Post("/students", handlers.CreateStudent)
	api.Put("/students/:roll_no", handlers.UpdateStudent)
	api.Delete("/students/:roll_no", handlers.DeleteStudent)

	return srv
}
