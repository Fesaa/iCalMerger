package log

import (
	"log/slog"
	"os"
)

type logger struct {
	*slog.Logger
	levelVar *slog.LevelVar
	notify   NotificationService
}

var Logger logger

func Init(level string, notify NotificationService) {
	Logger = logger{
		levelVar: new(slog.LevelVar),
		notify:   notify,
	}

	switch level {
	case "DEBUG":
		Logger.levelVar.Set(slog.LevelDebug)
	case "INFO":
		Logger.levelVar.Set(slog.LevelInfo)
	case "WARN":
		Logger.levelVar.Set(slog.LevelWarn)
	case "ERROR":
		Logger.levelVar.Set(slog.LevelError)
	default:
		Logger.levelVar.Set(slog.LevelInfo)
	}

	opts := slog.HandlerOptions{
		Level: Logger.levelVar,
	}
	Logger.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &opts))

	if notify.Url != "" {
		Logger.Info("Notifications enabled", "service", notify.Service)
	} else {
		Logger.Info("Notifications disabled")
	}
}

func (l *logger) Notify(msg string) {
	l.Debug("Sending notification", "message", msg)
	l.notify.process(msg)
}
