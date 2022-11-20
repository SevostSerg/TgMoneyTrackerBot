package configuration

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	BotToken     string
	DBFolderPath string
	DBFileName   string
}

const (
	configFileName string = "config.json"
)

var configurationInfo *Configuration

func GetInfo() *Configuration {
	if configurationInfo != nil {
		return configurationInfo
	}

	file, err := os.Open(configFileName)
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()
	err = json.NewDecoder(file).Decode(&configurationInfo)
	if err != nil {
		log.Panic(err)
	}

	return configurationInfo
}
