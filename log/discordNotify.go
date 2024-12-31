package log

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type DiscordNotification struct {
	Url string
}

func (n *DiscordNotification) Init(url string) {
	n.Url = url
}

func (n *DiscordNotification) Emit(msg string) {
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
