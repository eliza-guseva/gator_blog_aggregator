// Package config provides config for RSS feed aggregation.
package config

import (
	"encoding/json"
	"os"
	"fmt"
	"path/filepath"
)


type Config struct {
	DBUrl string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}


func Read() (*Config, error) {
	var config Config 
	home, err := os.UserHomeDir()
	if err != nil { return nil, fmt.Errorf("error getting home directory: %v", err) }
	configPath := filepath.Join(home, ".gatorconfig.json")
	
	data, err := os.ReadFile(configPath)
	if err != nil { return nil, fmt.Errorf("error reading config file: %v", err) }

	err = json.Unmarshal(data, &config)
	if err != nil { return nil, fmt.Errorf("error decoding config file: %v", err) }
	
	return &config, nil
}


func (c *Config) SetUser(user string) error {
	c.CurrentUser = user 
	home, err := os.UserHomeDir()
	if err != nil {return fmt.Errorf("error getting home directory: %v", err) }
	configPath := filepath.Join(home, ".gatorconfig.json")
	data, error := json.Marshal(c)
	if error != nil { return fmt.Errorf("error encoding config file: %v", error) }

	err = os.WriteFile(configPath, data, 0644)
	if err != nil { return fmt.Errorf("error writing config file: %v", err) }

	return nil
}
