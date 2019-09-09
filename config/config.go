package config

import (
	"encoding/json"
	"os"
)

// Config has some configuration parameters for the bot
type Config struct {
	DBCon        string `json:"db_con"`
	LVNLEndpoint string `json:"lvnl_endpoint"`
	BotKey       string `json:"bot_key"`
}

// Load loads the config file
func Load() *Config {
	// Load config
	file, err := os.Open("./config.json")
	if err != nil {
		panic("Coud not load config: " + err.Error())
	}
	defer file.Close()
	var config Config
	jsonParser := json.NewDecoder(file)
	jsonParser.Decode(config)

	return &config
}
