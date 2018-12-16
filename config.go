package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	WallpaperImgPath string `json:"wallpaperImgPath"`
	ClientConfPath   string `json:"clientConfPath"`
}

func configFromFile(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
