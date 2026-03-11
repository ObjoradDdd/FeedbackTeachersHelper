package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/ObjoradDdd/FeedbackTeachersHelper/docs"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/clients/llm"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/handlers"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/storage"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Feedback Teachers Helper API
// @version 1.0
// @description API для генерации фидбека учителям.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <your_token>
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("warning: no .env file found or error reading it")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "student"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "learn_orm"
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	dbStorage, err := storage.New(connString)
	if err != nil {
		log.Fatal("❌ Ошибка БД:", err)
	}

	if err := dbStorage.InitTables(); err != nil {
		log.Fatal("❌ Ошибка создания таблиц:", err)
	}

	defer dbStorage.Close()

	geminiClient := llm.NewGeminiClient()

	teacherService := services.NewTeacherService(dbStorage)
	groupService := services.NewGroupService(dbStorage)
	studentService := services.NewStudentService(dbStorage)
	tagService := services.NewTagService(dbStorage)
	feedbackService := services.NewFeedbackService(dbStorage, geminiClient)

	teacherHandler := handlers.NewTeacherHandler(teacherService)
	groupHandler := handlers.NewGroupHandler(groupService)
	studentHandler := handlers.NewStudentHandler(studentService)
	tagHandler := handlers.NewTagHandler(tagService)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)

	mux := http.NewServeMux()

	// Swagger UI
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	// Auth & Teacher
	mux.HandleFunc("POST /api/register", teacherHandler.Register)
	mux.HandleFunc("POST /api/login", teacherHandler.Login)
	mux.HandleFunc("POST /api/add_api_key", handlers.AuthMiddleware(teacherHandler.AddAPIKey))
	mux.HandleFunc("DELETE /api/delete_teacher", handlers.AuthMiddleware(teacherHandler.DeleteTeacher))

	// Groups
	mux.HandleFunc("POST /api/groups", handlers.AuthMiddleware(groupHandler.CreateGroup))
	mux.HandleFunc("GET /api/groups", handlers.AuthMiddleware(groupHandler.GetGroups))
	mux.HandleFunc("DELETE /api/groups/{id}", handlers.AuthMiddleware(groupHandler.DeleteGroup))
	mux.HandleFunc("PUT /api/groups/{id}", handlers.AuthMiddleware(groupHandler.UpdateGroup))

	// Students
	mux.HandleFunc("POST /api/students", handlers.AuthMiddleware(studentHandler.CreateStudent))
	mux.HandleFunc("GET /api/students/{groupId}", handlers.AuthMiddleware(studentHandler.GetStudentsGroup))
	mux.HandleFunc("DELETE /api/students/{id}", handlers.AuthMiddleware(studentHandler.DeleteStudent))
	mux.HandleFunc("PUT /api/students/{id}", handlers.AuthMiddleware(studentHandler.UpdateStudent))

	// Tags
	mux.HandleFunc("GET /api/tag", handlers.AuthMiddleware(tagHandler.GetTeacherTags))
	mux.HandleFunc("POST /api/tag", handlers.AuthMiddleware(tagHandler.CreateTag))
	mux.HandleFunc("PUT /api/tag/{id}", handlers.AuthMiddleware(tagHandler.UpdateTag))
	mux.HandleFunc("DELETE /api/tag/{id}", handlers.AuthMiddleware(tagHandler.DeleteTag))

	// Feedback
	mux.HandleFunc("POST /api/feedback", handlers.AuthMiddleware(feedbackHandler.GetFeedback))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
		slog.Info("User visited the site")
	})

	fmt.Println("🌐 Сервер запущен на https://localhost:443")
	if err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", handlers.EnableCORS(mux)); err != nil {
		log.Fatal("❌ Ошибка запуска сервера:", err)
	}
}
