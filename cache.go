package main

import (
	"encoding/json"
	"github.com/sachaos/todoist/lib"
	"io/ioutil"
	"os"
)

func LoadCache(filename string) (todoist.Sync, error) {
	sync, err := ReadCache(filename)
	if err != nil {
		err = WriteCache(default_cache_path, sync)
		if err != nil {
			return todoist.Sync{}, CommandFailed
		}
	}
	return sync, nil
}

func ReadCache(filename string) (todoist.Sync, error) {
	var s todoist.Sync
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return s, CommandFailed
	}
	err = json.Unmarshal(jsonString, &s)
	if err != nil {
		return s, CommandFailed
	}
	return s, nil
}

func WriteCache(filename string, sync todoist.Sync) error {
	buf, err := json.MarshalIndent(sync, "", "  ")
	if err != nil {
		return CommandFailed
	}
	err2 := ioutil.WriteFile(filename, buf, os.ModePerm)
	if err2 != nil {
		return CommandFailed
	}
	return nil
}
