package log

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type NotificationService struct {
	Service int
	Url     string
}

const (
	NotificationServiceTypeDiscord = iota
)

func (n *NotificationService) process(msg string) {
	if n.Url == "" {
		return
	}

	switch n.Service {
	case NotificationServiceTypeDiscord:
		n.toDiscord(msg)
	}
}

func (n *NotificationService) toDiscord(msg string) {
	payload := map[string]interface{}{
		"content":    msg,
		"username":   "iCal Merger Service",
		"avatar_url": "https://i.imgur.com/4M34hi2.png",
	}

	payloadJson, e := json.Marshal(payload)
	if e != nil {
		Logger.Error("Error marshaling webhook payload", "error", e)
		return
	}

	req, e := http.NewRequest("POST", n.Url, bytes.NewBuffer(payloadJson))
	if e != nil {
		Logger.Error("Error creating webhook request", "error", e)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	_, e = client.Do(req)
	if e != nil {
		Logger.Error("Error sending webhook", "error", e)
	}
}
