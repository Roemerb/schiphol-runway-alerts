package telegram

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
	"github.com/roemerb/schiphol-runway-alerts/config"
	"gopkg.in/gookit/color.v1"
)

// TelegramEndpoint is the URL where Telegram requests will be send to
var TelegramEndpoint = "https://api.telegram.org/"

// User represents a Telegram user
type User struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	ChatID       int
}

// Chat represents a chat (conversation) on Telegram
type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

// Message represents a Telegram message
type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

// SendMessage contains the data to send to Telegram to send a message
type SendMessage struct {
	ChatID     int    `url:"chat_id"`
	Text       string `url:"text"`
	ReplyToMsg int    `url:"reply_to_message_id"`
}

// Update is a new message from the Telegram webhook
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

// GetUpdateFromWebhook will convert the incoming JSON data from
// the Telegram webhook and transform it into an update
func GetUpdateFromWebhook(data []byte) (*Update, error) {
	var u = new(Update)
	err := json.Unmarshal(data, &u)
	if err != nil {
		return nil, err
	}

	return u, err
}

// Respond will digest an incoming message and produce a response
func Respond(update *Update) {
	color.Blue.Print(update.Message.From.FirstName + " (" + update.Message.From.Username + "): ")
	command := IdentifyCommand(update)
	if command != nil {
		color.Green.Print(update.Message.Text + "\n")
		command.Handle(update)
		return
	}
	color.Red.Print(update.Message.Text + "\n")

	msg := SendMessage{
		ChatID: update.Message.Chat.ID,
		Text:   "'" + update.Message.Text + "' is not a command",
	}
	msg.Send()
}

// Send will send a message
func (msg SendMessage) Send() error {
	q, _ := query.Values(msg)

	config := config.Load()
	url := TelegramEndpoint + config.BotKey + "/sendMessage?" + q.Encode()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(
			"Sending message failed with status " +
				response.Status + " (" + strconv.Itoa(response.StatusCode) + ")")
	}

	return nil
}
