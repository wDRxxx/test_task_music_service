package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/wDRxxx/test-task/internal/closer"
)

func SetupLogger(envLevel string, logsPath string) {
	var logger *slog.Logger

	switch envLevel {
	case "dev":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		t := time.Now().Format("02-01-06T15-04-05")
		logFilePath := fmt.Sprintf("%s/%s.log", logsPath, t)

		f, err := os.Create(logFilePath)
		if err != nil {
			log.Fatalf("error creating new log file: %v", err)
		}
		closer.Add(f.Close)

		logger = slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		}))
	}

	slog.SetDefault(logger)
}
