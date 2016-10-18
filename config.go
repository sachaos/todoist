package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func ParseConfig(filename string) (Config, interface{}) {
	var c Config
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, "Error, NotFound config file"
	}
	err = json.Unmarshal(jsonString, &c)
	if err != nil {
		return c, "Error, config file parse error"
	}
	return c, nil
}

func CreateConfig(filename string, config Config) interface{} {
	buf, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "Error, Failed to marshal json"
	}
	err2 := ioutil.WriteFile(filename, buf, os.ModePerm)
	if err2 != nil {
		return "Error, Failed to write config file"
	}
	return nil
}
