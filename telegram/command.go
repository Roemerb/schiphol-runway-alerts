package telegram

import (
	"strings"
)

// Command is a Telegram bot command
type Command interface {
	Command() string
	Handle(*Update) *SendMessage
}

// Commands is an array of commands
type Commands []Command

// ActiveCommands holds the commands that are currently active
var ActiveCommands = Commands{
	SubscribeCommand{
		CommandString: "subscribe",
	},
}

// SubscribeCommand is used to subscribe to the service over Telegram
type SubscribeCommand struct {
	CommandString string
}

// Command returns the name of the SubscibeCommand command
func (c SubscribeCommand) Command() string {
	return c.CommandString
}

// Handle handles responding to a SubscribeCommand
func (c SubscribeCommand) Handle(update *Update) *SendMessage {
	msg := SendMessage{
		ChatID:     update.Message.Chat.ID,
		Text:       update.Message.Text,
		ReplyToMsg: update.Message.MessageID,
	}
	msg.Send()

	return &msg
}

// IdentifyCommand will check if an incoming message is a command
func IdentifyCommand(update *Update) Command {
	text := update.Message.Text
	if len(text) < 2 || text[:1] != "/" {
		return nil
	}

	// Remove the /
	text = text[1:len(text)]

	text = strings.ToLower(text)
	words := strings.Split(text, " ")
	candidate := words[0]

	for _, c := range ActiveCommands {
		if c.Command() == candidate {
			return c
		}
	}

	return nil
}
