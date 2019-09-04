package main

import (
	"log"

	"github.com/roemerb/schiphol-runway-alerts/lvnl"
)

// VERSION is the version of the app
var VERSION = "0.1.0"

func main() {
	log.Print("Schiphol Runway Alerts version " + VERSION + "\n\nStarting...")

	changes := make(chan *lvnl.Runway)
	stop := make(chan bool)

	lvnl.Start(changes, stop)

	go func() {
		for {
			change := <-changes
			log.Println(change.Name + " has changed!")
			lvnl.PrintState()
		}
	}()
	<-stop
}
