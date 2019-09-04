package main

import (
	"log"
	"lvnl"
)

// VERSION is the version of the app
var VERSION = "0.1.0"

func main() {
	log.Print("Schiphol Runway Alerts version " + VERSION + "\n\nStarting...")

	changes := make(chan *lvnl.Runway)
	stop := make(chan bool)
}
