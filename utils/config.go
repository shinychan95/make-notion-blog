package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	DBPath  string `json:"db_path"`
	ApiKey  string `json:"api_key"`
	PostDir string `json:"post_directory"`
	ImgDir  string `json:"image_directory"`
	RootId  string `json:"root_id"`
}

func ReadConfig(configFilePath string) (Config, error) {
	var config Config

	configFile, err := os.Open(configFilePath)
	CheckError(err)
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	err = json.Unmarshal(byteValue, &config)

	return config, err
}
