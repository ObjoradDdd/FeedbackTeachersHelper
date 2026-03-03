package main

import (
	"log/slog"
	"os"

	"github.com/ObjoradDdd/FeedbackTeachersHelper/internal/storage"
)

func main() {
	slog.Info("Starting API")

	connString := "host=localhost port=5432 user=student password=password dbname=learn_orm sslmode=disable"

	db, err := storage.New(connString)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Storage successfully initialized!")

	if err := db.InitTables(); err != nil {
		slog.Error("Failed to initialize tables", "error", err)
		os.Exit(1)
	}
	slog.Info("Tables in database are ready!")

}
