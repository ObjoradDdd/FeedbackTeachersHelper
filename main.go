package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	_ "github.com/ObjoradDdd/FeedbackTeachersHelper/docs"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/clients/llm"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/handlers"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/kafka"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/services"
	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {

	// Создаем контекст, который будет отменен при получении сигнала завершения
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Создаем группу ожидания для управления горутинами
	wg := &sync.WaitGroup{}

	// Инициализируем логгер
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//  Канал для перехвата критических ошибок сервера
	serverErrors := make(chan error, 2)

	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	masterKey := os.Getenv("MASTER_KEY")
	if masterKey == "" {
		slog.Error("MASTER_KEY is not set in environment")
		return
	}

	// Инициализируем хранилище и запускаем сервер
	db, err := initStorage()
	if err != nil {
		slog.Error("failed to initialize storage", "error", err)
		return
	}
	defer db.Close()

	// Инициализируем сервисы
	gemini := llm.NewGeminiClient()
	userSvc := services.NewUserService(db, masterKey)
	groupSvc := services.NewGroupService(db)
	studentSvc := services.NewStudentService(db)
	tagSvc := services.NewTagService(db)
	fbSvc := services.NewFeedbackService(db, gemini, masterKey)

	// Запускаем Kafka consumer в отдельной горутине
	consumersCount, err := strconv.Atoi(os.Getenv("CONSUMERS_COUNT"))
	if err != nil {
		slog.Error("failed to parse consumers count", "error", err)
		return
	} else {
		startKafka(ctx, userSvc, wg, consumersCount, serverErrors)
	}

	// Настраиваем маршруты
	router := setupRouter(userSvc, groupSvc, studentSvc, tagSvc, fbSvc)

	// Запускаем HTTP сервер
	var srv *http.Server

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		slog.Warn("invalid SERVER_PORT", "error", err)
		return
	} else {
		srv = &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handlers.EnableCORS(router),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Info("server is starting", "port", port, "url", fmt.Sprintf("http://localhost:%d", port))
			serverErrors <- srv.ListenAndServe()
		}()

	}

	// Ожидаем сигнала завершения или ошибки запуска приложения
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			slog.Error("server failed prematurely", "error", err)
			stop()
		}
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	// Создаем контекст с таймаутом для корректного завершения сервера
	if srv != nil {
		shutdownCtx, serverStop := context.WithTimeout(context.Background(), 10*time.Second)
		defer serverStop()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("forced shutdown", "error", err)
			srv.Close()
		}
	}

	// Ждем завершения всех горутин (HTTP сервер и Kafka consumer)
	wg.Wait()
	slog.Info("shutting down server")
}

func initStorage() (*storage.Storage, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := storage.New(conn)
	if err != nil {
		return nil, fmt.Errorf("db connection: %w", err)
	}

	return db, nil
}

func setupRouter(userSvc *services.UserService, groupSvc *services.GroupService, studentSvc *services.StudentService, tagSvc *services.TagService, fbSvc *services.FeedbackService) *http.ServeMux {
	mux := http.NewServeMux()

	validator := validator.New()

	hUser := handlers.NewUserHandler(userSvc, validator)
	hGroup := handlers.NewGroupHandler(groupSvc, validator)
	hStudent := handlers.NewStudentHandler(studentSvc, validator)
	hTag := handlers.NewTagHandler(tagSvc, validator)
	hFb := handlers.NewFeedbackHandler(fbSvc, validator)

	// Роуты
	mux.HandleFunc("POST /api/add_api_key", handlers.AuthMiddleware(hUser.AddAPIKey))
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

func startKafka(ctx context.Context, userSvc *services.UserService, wg *sync.WaitGroup, consumersCount int, serverErrors chan<- error) {
	manager := kafka.NewConsumerManager(
		os.Getenv("KAFKA_BROKERS"),
		userSvc,
		"fth_group",
		"user_events",
		consumersCount,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := manager.Start(ctx); err != nil {
			slog.Error("Kafka manager stopped with error", "error", err)
			serverErrors <- err
		}
	}()
}
