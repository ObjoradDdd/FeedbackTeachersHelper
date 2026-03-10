package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

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

	dbStorage.InitTables()

	defer dbStorage.Close()

	teacherService := services.NewTeacherService(dbStorage)
	groupService := services.NewGroupService(dbStorage)
	studentService := services.NewStudentService(dbStorage)
	tagService := services.NewTagService(dbStorage)

	teacherHandler := handlers.NewTeacherHandler(teacherService)
	groupHandler := handlers.NewGroupHandler(groupService)
	studentHandler := handlers.NewStudentHandler(studentService)
	tagHandler := handlers.NewTagHandler(tagService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", teacherHandler.Register)
	mux.HandleFunc("POST /api/login", teacherHandler.Login)
	mux.HandleFunc("POST /api/add_api_key", handlers.AuthMiddleware(teacherHandler.AddAPIKey))
	mux.HandleFunc("DELETE /api/delete_teacher", handlers.AuthMiddleware(teacherHandler.DeleteTeacher))

	mux.HandleFunc("POST /api/groups", handlers.AuthMiddleware(groupHandler.CreateGroup))
	mux.HandleFunc("GET /api/groups", handlers.AuthMiddleware(groupHandler.GetGroups))
	mux.HandleFunc("DELETE /api/groups", handlers.AuthMiddleware(groupHandler.DeleteGroup))
	mux.HandleFunc("PUT /api/groups", handlers.AuthMiddleware(groupHandler.UpdateGroup))

	mux.HandleFunc("POST /api/students", handlers.AuthMiddleware(studentHandler.CreateStudent))
	mux.HandleFunc("GET /api/students", handlers.AuthMiddleware(studentHandler.GetStudentsGroup))
	mux.HandleFunc("DELETE /api/students", handlers.AuthMiddleware(studentHandler.DeleteStudent))
	mux.HandleFunc("PUT /api/students", handlers.AuthMiddleware(studentHandler.UpdateStudent))

	mux.HandleFunc("GET /api/tag", handlers.AuthMiddleware(tagHandler.GetTeacherTags))
	mux.HandleFunc("POST /api/tag", handlers.AuthMiddleware(tagHandler.CreateTag))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
		slog.Info("User visited the site")
	})

	fmt.Println("🌐 Сервер запущен на https://localhost:8080")
	if err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", mux); err != nil {
		log.Fatal("❌ Ошибка запуска сервера:", err)
	}
}
