package lib

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SlackParams struct {
	Text string `json:"text"`
}

func SendToSlack(webhook string, message string) bool {
	bjsonStr, _ := json.Marshal(SlackParams{Text: message})

	r, _ := http.NewRequest("POST", webhook, bytes.NewBuffer(bjsonStr))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	defer resp.Body.Close()

	return true
}
