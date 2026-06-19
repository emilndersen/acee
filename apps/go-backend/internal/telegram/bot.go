package telegram

import (
	"fmt"
	"net/http"
	"net/url"
)

type Bot struct {
	token  string
	chatID string
}

func NewBot(token, chatID string) *Bot {
	return &Bot{token: token, chatID: chatID}
}

func (b *Bot) Enabled() bool {
	return b.token != "" && b.chatID != ""
}

func (b *Bot) Send(text string) error {
	if !b.Enabled() {
		return nil
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id":    {b.chatID},
		"text":       {text},
		"parse_mode": {"HTML"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
