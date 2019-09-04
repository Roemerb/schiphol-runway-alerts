package http

import (
	"io/ioutil"
	"net/http"

	"github.com/roemerb/schiphol-runway-alerts/telegram"
)

// HandleTelegramWebhook is a http handler func for the Telegram webhook
func HandleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	json, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(422)
		panic(err)
	}

	update, err := telegram.GetUpdateFromWebhook(json)
	if err != nil {
		w.WriteHeader(422)
	}

	w.WriteHeader(http.StatusOK)

	telegram.Respond(update)
	defer r.Body.Close()
}
