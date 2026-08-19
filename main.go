package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  *struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

type GetUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func main() {
	token := os.Getenv("TELE_TOKEN")
	if token == "" {
		log.Println("TELE_TOKEN is not set; idling instead of crashing")
		select {}
	}

	base := "https://api.telegram.org/bot" + token
	client := &http.Client{Timeout: 60 * time.Second}
	offset := 0

	log.Println("bot-project started, polling Telegram API")

	for {
		reqURL := fmt.Sprintf("%s/getUpdates?timeout=50&offset=%d", base, offset)
		resp, err := client.Get(reqURL)
		if err != nil {
			log.Printf("getUpdates error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var upd GetUpdatesResponse
		if err := json.Unmarshal(body, &upd); err != nil {
			log.Printf("decode error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if !upd.OK {
			log.Printf("telegram api error: %s", string(body))
			time.Sleep(5 * time.Second)
			continue
		}

		for _, u := range upd.Result {
			offset = u.UpdateID + 1
			if u.Message == nil {
				continue
			}
			form := url.Values{}
			form.Set("chat_id", fmt.Sprintf("%d", u.Message.Chat.ID))
			form.Set("text", "echo: "+u.Message.Text)
			if _, err := client.PostForm(base+"/sendMessage", form); err != nil {
				log.Printf("sendMessage error: %v", err)
			}
		}
	}
}
