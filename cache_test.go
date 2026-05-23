package main

import (
	"encoding/json"
	"os"
	"testing"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempCache(t *testing.T, content map[string]interface{}) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cache-*.json")
	require.NoError(t, err)
	defer f.Close()
	require.NoError(t, json.NewEncoder(f).Encode(content))
	return f.Name()
}

func TestReadCache_FileNotFound(t *testing.T) {
	var s todoist.Store
	err := ReadCache("/nonexistent/path/cache.json", &s)
	assert.Error(t, err)
}

func assertCacheFileSchemaVersion(t *testing.T, path string, want int) {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var stored struct {
		SchemaVersion int `json:"schema_version"`
	}
	require.NoError(t, json.Unmarshal(data, &stored))
	assert.Equal(t, want, stored.SchemaVersion, "schema_version in cache file on disk")
}

func TestReadCache_OldCache_NoSchemaVersion(t *testing.T) {
	// Old caches have no schema_version field; json.Unmarshal sets missing
	// int fields to 0, which != currentSchemaVersion, so SyncToken must be
	// reset. This test verifies both the Go behavior and the reset logic.
	path := writeTempCache(t, map[string]interface{}{
		"sync_token": "abc123",
		"items":      []interface{}{},
	})
	var s todoist.Store
	err := ReadCache(path, &s)
	assert.NoError(t, err)
	assert.Equal(t, "*", s.SyncToken)
	assert.Equal(t, currentSchemaVersion, s.SchemaVersion)
	assertCacheFileSchemaVersion(t, path, currentSchemaVersion)
}

func TestReadCache_WrongSchemaVersion(t *testing.T) {
	path := writeTempCache(t, map[string]interface{}{
		"sync_token":     "abc123",
		"schema_version": 99,
		"items":          []interface{}{},
	})
	var s todoist.Store
	err := ReadCache(path, &s)
	assert.NoError(t, err)
	assert.Equal(t, "*", s.SyncToken)
	assert.Equal(t, currentSchemaVersion, s.SchemaVersion)
	assertCacheFileSchemaVersion(t, path, currentSchemaVersion)
}

func TestReadCache_CorrectSchemaVersion_PreservesSyncToken(t *testing.T) {
	path := writeTempCache(t, map[string]interface{}{
		"sync_token":     "mytoken",
		"schema_version": 1,
		"items":          []interface{}{},
		"projects":       []interface{}{},
		"labels":         []interface{}{},
		"sections":       []interface{}{},
	})
	var s todoist.Store
	err := ReadCache(path, &s)
	assert.NoError(t, err)
	assert.Equal(t, "mytoken", s.SyncToken)
	assert.Equal(t, 1, s.SchemaVersion)
}
