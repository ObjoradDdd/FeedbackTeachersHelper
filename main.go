package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ObjoradDdd/FeedbackTeachersHelper/docs"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/clients/llm"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/handlers"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/kafka"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	db, err := initStorage()
	if err != nil {
		slog.Error("failed to initialize storage", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	router := setupRouter(db)
	startKafka(db)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "5135"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handlers.EnableCORS(router),
	}

	go func() {
		slog.Info("server is starting", "port", port, "url", "http://localhost:"+port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("received shutdown signal", "signal", (<-quit).String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}
	slog.Info("server exited gracefully")
}

func initStorage() (*storage.Storage, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := storage.New(conn)
	if err != nil {
		return nil, fmt.Errorf("db connection: %w", err)
	}

	if err := db.InitTables(); err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}
	return db, nil
}

func setupRouter(db *storage.Storage) *http.ServeMux {
	mux := http.NewServeMux()

	gemini := llm.NewGeminiClient()
	userSvc := services.NewUserService(db)
	groupSvc := services.NewGroupService(db)
	studentSvc := services.NewStudentService(db)
	tagSvc := services.NewTagService(db)
	fbSvc := services.NewFeedbackService(db, gemini)

	hUser := handlers.NewUserHandler(userSvc)
	hGroup := handlers.NewGroupHandler(groupSvc)
	hStudent := handlers.NewStudentHandler(studentSvc)
	hTag := handlers.NewTagHandler(tagSvc)
	hFb := handlers.NewFeedbackHandler(fbSvc)

	// Роуты остаются прежними
	mux.HandleFunc("POST /api/add_api_key", handlers.AuthMiddleware(hUser.AddAPIKey))
	mux.HandleFunc("DELETE /api/delete_user", handlers.AuthMiddleware(hUser.DeleteUser))
	mux.HandleFunc("POST /api/groups", handlers.AuthMiddleware(hGroup.CreateGroup))
	mux.HandleFunc("GET /api/groups", handlers.AuthMiddleware(hGroup.GetGroups))
	mux.HandleFunc("PUT /api/groups/{id}", handlers.AuthMiddleware(hGroup.UpdateGroup))
	mux.HandleFunc("DELETE /api/groups/{id}", handlers.AuthMiddleware(hGroup.DeleteGroup))
	mux.HandleFunc("POST /api/students", handlers.AuthMiddleware(hStudent.CreateStudent))
	mux.HandleFunc("GET /api/students/{groupId}", handlers.AuthMiddleware(hStudent.GetStudentsGroup))
	mux.HandleFunc("PUT /api/students/{id}", handlers.AuthMiddleware(hStudent.UpdateStudent))
	mux.HandleFunc("DELETE /api/students/{id}", handlers.AuthMiddleware(hStudent.DeleteStudent))
	mux.HandleFunc("GET /api/tag", handlers.AuthMiddleware(hTag.GetUserTags))
	mux.HandleFunc("POST /api/tag", handlers.AuthMiddleware(hTag.CreateTag))
	mux.HandleFunc("PUT /api/tag/{id}", handlers.AuthMiddleware(hTag.UpdateTag))
	mux.HandleFunc("DELETE /api/tag/{id}", handlers.AuthMiddleware(hTag.DeleteTag))
	mux.HandleFunc("POST /api/feedback", handlers.AuthMiddleware(hFb.GetFeedback))

	return mux
}

func startKafka(db *storage.Storage) {
	userSvc := services.NewUserService(db)
	consumer := kafka.NewUserConsumer(
		os.Getenv("KAFKA_BROKERS"),
		userSvc,
		"fth_group",
		"user_events",
	)
	go consumer.Start()
}
