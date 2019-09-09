package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route holds all the info for a route
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// Routes is just a Route array
type Routes []Route

var routes = Routes{
	Route{
		"TelegramWebhook",
		http.MethodPost,
		"/telegram/webhook",
		HandleTelegramWebhook,
	},
}

// NewRouter generates a new mux.Router instance
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.Handler

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
