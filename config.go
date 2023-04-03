package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	DBPath          string `json:"db_path"`
	OutputDirectory string `json:"output_directory"`
}

func readConfig(configFilePath string) (Config, error) {
	var config Config

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	err = json.Unmarshal(byteValue, &config)

	return config, err
}
