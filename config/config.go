package config

import (
	"encoding/json"
	"io/ioutil"
)

const configFileName = "config.json"

type Config struct {
	GetTeamURL           string
	TargetTeams          []string
	TimeoutInMilliSecond int
}

// GetTargetTeams returns a map of target teams
func (c *Config) GetTargetTeams() map[string]bool {
	targetTeams := make(map[string]bool)
	for _, team := range c.TargetTeams {
		targetTeams[team] = true
	}

	return targetTeams
}

// Get reads the configuration from the config file
func Get() (*Config, error) {
	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
