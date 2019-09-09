package lvnl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/roemerb/schiphol-runway-alerts/config"
)

// Runways is a definition table for the different runways at Schiphol Airport
var Runways = map[string]string{
	"18L": "Aalsmeerbaan",
	"36R": "Aalsmeerbaan",
	"09":  "Buitenveldertbaan",
	"27":  "Buitenveldertbaan",
	"06":  "Kaagbaan",
	"24":  "Kaagbaan",
	"18R": "Polderbaan",
	"36L": "Polderbaan",
	"18C": "Zwanenburgbaan",
	"36C": "Zwanenburgbaan",
}

// RunwayUsageRequest contains the parameters needed to perform
// a request to the LVNL service to retrieve runway usage information
type RunwayUsageRequest struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
}

// RunwayUsageResponse contains the response from a runway usage
// request. It's the raw unmarshalled JSON
type RunwayUsageResponse struct {
	ID       int       `json:"Id"`
	Updated  time.Time `json:"Updated"`
	Start    time.Time `json:"Start"`
	End      time.Time `json:"End"`
	Landing1 string    `json:"Landing1"`
	Landing2 string    `json:"Landing2"`
	Landing3 string    `json:"Landing3"`
	Takeoff1 string    `json:"Takeoff1"`
	Takeoff2 string    `json:"Takeoff2"`
	Takeoff3 string    `json:"Takeoff3"`
	State    string    `json:"State"`
	IsLast   bool      `json:"isLast"`
}

// GetActiveLandingRunways takes the response from the LVNL api and extracts
// the currently active runways, as they're all in seperate fields. This makes
// for some cleaner code while parsing
func (res RunwayUsageResponse) GetActiveLandingRunways() []string {
	var active []string
	if res.Landing1 != "" {
		active = append(active, res.Landing1)
	}

	if res.Landing2 != "" {
		active = append(active, res.Landing2)
	}

	if res.Landing3 != "" {
		active = append(active, res.Landing3)
	}

	return active
}

// GetActiveTakeoffRunways takes the response from the LVNL api and extracts
// the currently active runways, as they're all in seperate fields. This makes
// for some cleaner code while parsing
func (res RunwayUsageResponse) GetActiveTakeoffRunways() []string {
	var active []string
	if res.Takeoff1 != "" {
		active = append(active, res.Takeoff1)
	}

	if res.Takeoff2 != "" {
		active = append(active, res.Takeoff2)
	}

	if res.Takeoff3 != "" {
		active = append(active, res.Takeoff3)
	}

	return active
}

// GetRunwayUsage uses a RunwayUsageRequest to fetch the current runway usage
// from LVNL
func GetRunwayUsage(req *RunwayUsageRequest) RunwayUsageResponse {
	payloadArr := []string{
		strconv.Itoa(req.Year),
		strconv.Itoa(req.Month),
		strconv.Itoa(req.Day),
		strconv.Itoa(req.Hour),
		strconv.Itoa(req.Minute),
	}
	payload := "[" + strings.Join(payloadArr, ",") + "]"

	config := config.Load()
	resp, err := http.Post(config.LVNLEndpoint, "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal("FAILED: " + err.Error())
	}

	var response RunwayUsageResponse
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &response)

	return response
}
