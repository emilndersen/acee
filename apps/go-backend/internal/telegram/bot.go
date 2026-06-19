package telegram

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Bot struct {
	token  string
	chatID string
}

func NewBot(token, chatID string) *Bot {
	log.Printf("Telegram bot init: token=%q, chatID=%q", token, chatID)
	return &Bot{token: token, chatID: chatID}
}

func (b *Bot) Enabled() bool {
	return b.token != "" && b.chatID != ""
}

func (b *Bot) Send(text string) error {
	if !b.Enabled() {
		log.Println("Telegram bot disabled: token or chatID empty")
		return nil
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id":    {b.chatID},
		"text":       {text},
		"parse_mode": {"HTML"},
	})
	if err != nil {
		log.Printf("Telegram send error: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("Telegram send status: %d", resp.StatusCode)
	return nil
}

func (b *Bot) SendBookingNotification(name, contact, shootType, date, idea string) error {
	text := fmt.Sprintf(
		"📸 <b>Новая заявка на съёмку!</b>\n\n"+
			"<b>Имя:</b> %s\n"+
			"<b>Контакт:</b> %s\n"+
			"<b>Тип:</b> %s\n"+
			"<b>Дата:</b> %s\n"+
			"<b>Идея:</b> %s",
		name, contact, shootType, date, idea,
	)

	if idea == "" {
		text = fmt.Sprintf(
			"📸 <b>Новая заявка на съёмку!</b>\n\n"+
				"<b>Имя:</b> %s\n"+
				"<b>Контакт:</b> %s\n"+
				"<b>Тип:</b> %s\n"+
				"<b>Дата:</b> %s",
			name, contact, shootType, date,
		)
	}

	return b.Send(text)
}
