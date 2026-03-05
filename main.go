package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/handlers"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/storage"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("warning: no .env file found or error reading it")
	}

	connString := "host=localhost port=5432 user=student password=password dbname=learn_orm sslmode=disable"
	dbStorage, err := storage.New(connString)
	if err != nil {
		log.Fatal("❌ Ошибка БД:", err)
	}
	defer dbStorage.Close()

	teacherService := services.NewTeacherService(dbStorage)
	groupService := services.NewGroupService(dbStorage)

	teacherHandler := handlers.NewTeacherHandler(teacherService)
	groupHandler := handlers.NewGroupHandler(groupService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", teacherHandler.Register)
	mux.HandleFunc("POST /api/login", teacherHandler.Login)
	mux.HandleFunc("POST /api/add_api_key", handlers.AuthMiddleware(teacherHandler.AddAPIKey))
	mux.HandleFunc("DELETE /api/delete_teacher", handlers.AuthMiddleware(teacherHandler.DeleteTeacher))

	mux.HandleFunc("POST /api/groups", handlers.AuthMiddleware(groupHandler.CreateGroup))
	mux.HandleFunc("GET /api/groups", handlers.AuthMiddleware(groupHandler.GetGroups))
	mux.HandleFunc("DELETE /api/groups", handlers.AuthMiddleware(groupHandler.DeleteGroup))
	mux.HandleFunc("PUT /api/groups", handlers.AuthMiddleware(groupHandler.UpdateGroup))

	fmt.Println("🌐 Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("❌ Ошибка запуска сервера:", err)
	}
}
