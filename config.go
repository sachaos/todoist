package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func LoadConfig(filename string) (Config, error) {
	var err error
	var token string
	config := Config{}
	config, err = ParseConfigFile(filename)
	if err != nil {
		fmt.Printf("Input API Token: ")
		fmt.Scan(&token)
		config = Config{Token: token}
		err = CreateConfigFile(filename, config)
		if err != nil {
			return config, CommandFailed
		}
	}
	return config, nil
}

func ParseConfigFile(filename string) (Config, error) {
	var c Config
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, CommandFailed
	}
	err = json.Unmarshal(jsonString, &c)
	if err != nil {
		return c, CommandFailed
	}
	return c, nil
}

func CreateConfigFile(filename string, config Config) error {
	buf, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return CommandFailed
	}
	err2 := ioutil.WriteFile(filename, buf, os.ModePerm)
	if err2 != nil {
		return CommandFailed
	}
	return nil
}
