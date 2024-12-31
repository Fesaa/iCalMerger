package log

const (
	NotificationServiceTypeDiscord = "discord"
)

func NewNotificationService(service string, url string) NotificationService {
	var s NotificationService

	switch service {
	case NotificationServiceTypeDiscord:
		s = &DiscordNotification{}
	default:
		Logger.Info("Notifications disabled")
		return nil
	}

	Logger.Info("Notifications enabled", "service", service, "url", url)
	s.Init(url)
	return s
}

type NotificationService interface {
	Init(url string)
	Emit(msg string)
}
