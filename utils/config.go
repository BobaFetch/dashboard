package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	ADDRES       string `json:"ADDRESS"`
	PORT         int    `json:"PORT"`
	SQL_HOST     string `json:"SQL_HOST"`
	SQL_USER     string `json:"SQL_USER"`
	SQL_PASSWORD string `json:"SQL_PASSWORD"`
	SQL_DB       string `json:"SQL_DB"`
	SQL_PORT     int    `json:"SQL_PORT"`
}

var config Config

func LoadConfig() (*Config, error) {
	file, err := os.Open("utils/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
