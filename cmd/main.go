package main

import (
	"log/slog"

	"github.com/asullivan219/receiptProcessor/internal/server"
)

func main() {
	server := server.NewServer(8080, "resources/database.db")

	slog.Info("Starting server",
		"port", 8080,
		"db file", "resources/database.db",
	)
	slog.Error(
		"Error starting server",
		"error", server.Srv.ListenAndServe(),
	)

}
