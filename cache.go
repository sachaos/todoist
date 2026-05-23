package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/sachaos/todoist/lib"
)

const currentSchemaVersion = 1

func LoadCache(filename string, s *todoist.Store) error {
	err := ReadCache(filename, s)
	if err != nil {
		err = WriteCache(cachePath, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadCache(filename string, s *todoist.Store) error {
	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return CommandFailed
	}

	// Two-pass: check schema version before full unmarshal so that a
	// schema change never leaves the cache in a broken state.
	var meta struct {
		SchemaVersion int `json:"schema_version"`
	}
	json.Unmarshal(jsonBytes, &meta) // error ignored: missing field yields 0

	if meta.SchemaVersion != currentSchemaVersion {
		// Old or mismatched cache: force a full resync on next sync call.
		s.SyncToken = "*"
		s.SchemaVersion = currentSchemaVersion
		_ = WriteCache(filename, s)
		return nil
	}

	if err := json.Unmarshal(jsonBytes, s); err != nil {
		return CommandFailed
	}
	s.ConstructItemTree()
	return nil
}

func WriteCache(filename string, s *todoist.Store) error {
	buf, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return CommandFailed
	}
	err = AssureExists(filename)
	if err != nil {
		return err
	}
	err2 := ioutil.WriteFile(filename, buf, os.ModePerm)
	if err2 != nil {
		return errors.New("Couldn't write to the cache file")
	}
	return nil
}
