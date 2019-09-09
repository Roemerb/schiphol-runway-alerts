package main

import (
	"log"
	sysHttp "net/http"
	"strconv"

	"github.com/roemerb/schiphol-runway-alerts/http"
	"github.com/roemerb/schiphol-runway-alerts/lvnl"
	"github.com/roemerb/schiphol-runway-alerts/telegram"
	"gopkg.in/gookit/color.v1"
)

// VERSION is the version of the app
var VERSION = "0.1.0"

// Config has some configuration parameters for the bot
type Config struct {
	DBCon        string `json:"db_con"`
	LVNLEndpoint string `json:"lvnl_endpoint"`
	BotKey       string `json:"bot_key"`
}

func main() {
	log.Print("Schiphol Runway Alerts version " + VERSION + "\n\nStarting...")

	changes := make(chan *lvnl.Runway)
	stop := make(chan bool)

	// Start listening for runway changes
	lvnl.Start(changes, stop)

	// When a change occurs, notify all subscribers
	go func() {
		for {
			change := <-changes
			log.Println(change.Name + " has changed!")
			lvnl.PrintState()

			var repo telegram.SubscriberRepository
			subs, err := repo.GetAllSubscribers()
			if err != nil {
				log.Println(err.Error())
				break
			}
			color.Green.Println("Notifying " + strconv.Itoa(len(subs)) + " of change")
			for i, sub := range subs {
				color.Green.Print(strconv.Itoa(i) + "/" + strconv.Itoa(len(subs)) + sub.Username + ": ")
				msg := telegram.SendMessage{
					ChatID: sub.TelegramChatID,
					Text:   change.ToMessage(),
				}
				err := msg.Send()
				if err != nil {
					color.Red.Print("FAILED: " + err.Error() + "\n")
				} else {
					color.Green.Print("SUCCESS\n")
				}
			}
		}
	}()

	// Start HTTP server to listen for incoming Telegram messages
	sysHttp.ListenAndServe(":3000", http.NewRouter())
}
