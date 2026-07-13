package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New() *slog.Logger {
	var Level = new(slog.LevelVar)
	// مقداردهی سطح لاگ از روی محیط برنامه (.env یا سیستم‌عامل)
	Level.Set(parseLevel(os.Getenv("LOG_LEVEL")))

	// ایجاد هندر JSON پیش‌فرض برای خروجی استاندارد
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: Level,
	})

	// تزریق اطلاعات اولیه سرویس به تمام لاگ‌ها
	return slog.New(handler).With(
		"service", os.Getenv("APP_NAME"),
		"version", os.Getenv("APP_VERSION"),
	)
}

func parseLevel(v string) slog.Level {
	switch strings.ToLower(v) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
