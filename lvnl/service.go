package lvnl

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/gookit/color.v1"
)

// Service describes the service that will listen for runway usage changes
type Service interface {
	Start(chan *Runway, chan bool)
}

// Runway contains the naming and state information of a runway
type Runway struct {
	Code      string
	Name      string
	Direction bool // 0 for takeoff, 1 for landing
	Active    bool
}

// State contains the current activity state of all runways
var State map[string]Runway

// IsActive indicates if rw is currently active
func (rw Runway) IsActive() bool {
	if rwState, ok := State[rw.Code]; ok {
		return rwState.Active
	}

	return false
}

// Start starts the LVNL service that will watch runway activity from
// the API. When a change is detected, the updated Runway will be pushed
// onto the channel that this function returns
func Start(changes chan *Runway, stop chan bool) {
	log.Println("Starting LVNL service")
	initiateState()
	t := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-stop:
				log.Println("LVNL Service received stop signal. Killing service")
				return
			case <-t.C:
				log.Println("Fetching runway usage")
				runwayUsageRequest := RunwayUsageRequest{
					Year:   time.Now().Year(),
					Month:  int(time.Now().Month()),
					Day:    time.Now().Day(),
					Hour:   time.Now().Hour(),
					Minute: time.Now().Minute(),
				}

				response := GetRunwayUsage(&runwayUsageRequest)
				updateState(&response, changes)
			}
		}
	}()
}

// PrintState will pretty-print the current state to the console
func PrintState() {
	longestLength := 0
	for _, rw := range State {
		full := rw.Name + " (" + rw.Code + ")"
		if len(full) > longestLength {
			longestLength = len(full)
		}
	}

	for _, rw := range State {
		full := rw.Name + " (" + rw.Code + ")"
		fmt.Print(full)
		if len(full) < longestLength {
			for i := 0; i < longestLength-len(full); i++ {
				fmt.Print(" ")
			}
		}
		fmt.Print("\t\t\t")

		if rw.Active {
			color.Green.Print("ACTIVE\t\t")
			if rw.Direction {
				color.Yellow.Print("Landing\n")
			} else {
				color.Yellow.Print("Takeoff\n")
			}
		} else {
			color.Red.Print("INACTIVE\n")
		}
	}
}

func initiateState() {
	State = make(map[string]Runway)
	for code, name := range Runways {
		rw := Runway{
			Code:      code,
			Name:      name,
			Direction: false,
			Active:    false,
		}
		State[code] = rw
	}
	log.Println("State initiated!")
}

func updateState(res *RunwayUsageResponse, changes chan *Runway) {
	// Get the active takeoff and landing runways
	activeTakeoff := res.GetActiveTakeoffRunways()
	activeLanding := res.GetActiveLandingRunways()
	allActive := append([]string{}, append(activeTakeoff, activeLanding...)...)

	// Iterate over the current state to detect updates
	for code, tw := range State {
		active := strArrContains(code, allActive)
		if active { // The runway is currently active
			if !tw.Active {
				// If not update the state and push change into channel
				tw.Active = true
				tw.Direction = strArrContains(code, activeLanding)
				State[code] = tw
				changes <- &tw
			}
		} else {
			// Also update if the runway is no longer active, but is
			// marked as active in the state
			if tw.Active {
				tw.Active = false
				tw.Direction = false
				State[code] = tw
				changes <- &tw
			}
		}
	}
}

func strArrContains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}
