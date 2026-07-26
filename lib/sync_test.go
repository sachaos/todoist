package todoist

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Regression test: a full sync must fully replace the cached store. Decoding
// the response into an already-populated store merges each incoming item into
// whatever cached item previously occupied the same slice index, and a JSON
// null section_id leaves that unrelated item's stale SectionID in place.
func TestSyncReplacesCachedStore(t *testing.T) {
	response := `{
		"full_sync": true,
		"sync_token": "token-2",
		"items": [
			{"id": "2", "content": "task two", "project_id": "p1", "section_id": null},
			{"id": "1", "content": "task one", "project_id": "p1", "section_id": "s1"}
		],
		"projects": [{"id": "p1", "name": "Project"}],
		"labels": []
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer server.Close()

	origServer := Server
	Server = server.URL
	defer func() { Server = origServer }()

	client := NewClient(&Config{AccessToken: "dummy"})

	// Simulate the store loaded from cache.json: same tasks, but ordered
	// differently than the sync response, with item 1 holding a section.
	cached := &Store{
		SyncToken: "token-1",
		Items: Items{
			{BaseItem: BaseItem{HaveID: HaveID{ID: "1"}, Content: "task one", HaveProjectID: HaveProjectID{ProjectID: "p1"}}, HaveSectionID: HaveSectionID{SectionID: "s1"}},
			{BaseItem: BaseItem{HaveID: HaveID{ID: "2"}, Content: "task two", HaveProjectID: HaveProjectID{ProjectID: "p1"}}},
		},
		Projects: Projects{
			{HaveID: HaveID{ID: "p1"}, Name: "Project"},
		},
	}
	cached.ConstructItemTree()
	client.Store = cached

	err := client.Sync(context.Background())
	assert.NoError(t, err)

	itemOne := client.Store.FindItem("1")
	itemTwo := client.Store.FindItem("2")
	assert.NotNil(t, itemOne)
	assert.NotNil(t, itemTwo)

	assert.Equal(t, "s1", itemOne.SectionID, "item 1 should keep its section")
	assert.Equal(t, "", itemTwo.SectionID, "item 2 has a null section_id and must not inherit a stale section from the cache")
	assert.Equal(t, "token-2", client.Store.SyncToken)
}
