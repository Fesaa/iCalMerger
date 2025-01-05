package log

import (
	"log/slog"
	"os"

	"github.com/Fesaa/ical-merger/config"
)

type logger struct {
	*slog.Logger
	levelVar *slog.LevelVar
	notify   NotificationService
}

var Logger logger

func Init(level string, notify config.Notification) {
	Logger = logger{
		levelVar: new(slog.LevelVar),
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

	Logger.notify = NewNotificationService(notify.Service, notify.Url)
}

func (l *logger) Notify(msg string) {
	if l.notify == nil {
		return
	}

	l.Debug("Sending notification", "message", msg)
	l.notify.Emit(msg)
}
