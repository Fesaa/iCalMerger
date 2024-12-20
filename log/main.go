package log

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/withmandala/go-log"
)

var Log *log.Logger

func Init(debug bool) {
	Log = log.New(os.Stderr)
	if debug {
		Log.WithDebug()
	}
}

func ToWebhook(url string, msg string) {
	if url == "" {
		return
	}

	Log.Debug("Sending webhook with content: ", msg)
	go func() {
		paylout := map[string]interface{}{
			"content":    msg,
			"username":   "iCal Merger Service",
			"avatar_url": "https://i.imgur.com/4M34hi2.png",
		}

		payloadJson, e := json.Marshal(paylout)
		if e != nil {
			Log.Error("Error marshaling webhook payload", e)
			return
		}

		req, e := http.NewRequest("POST", url, bytes.NewBuffer(payloadJson))
		if e != nil {
			Log.Error("Error creating webhook request", e)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		_, e = client.Do(req)
		if e != nil {
			Log.Error("Error sending webhook", e)
		}
	}()
}
