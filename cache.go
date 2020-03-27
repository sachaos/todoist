package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"errors"
	"path/filepath"

	"github.com/sachaos/todoist/lib"
)

func LoadCache(filename string, s *todoist.Store) error {
	err := ReadCache(filename, s)
	if err != nil {
		err = WriteCache(default_cache_path, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadCache(filename string, s *todoist.Store) error {
	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return CommandFailed
	}
	err = json.Unmarshal(jsonString, &s)
	if err != nil {
		return CommandFailed
	}
	s.ConstructItemTree()
	return nil
}

func createCache(filename string) error {
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't create cache file '%s'", filename))
	}
	return nil
}

func WriteCache(filename string, s *todoist.Store) error {
	buf, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return CommandFailed
	}
	_, fileErr := os.Stat(filename)
	if os.IsNotExist(fileErr) {
		err = createCache(filename)
		if err != nil {
			return err
		}
	}
	err2 := ioutil.WriteFile(filename, buf, os.ModePerm)
	if err2 != nil {
		return errors.New("Couldn't write to the cache file")
	}
	return nil
}
