package todoist

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestBaseURL_Default(t *testing.T) {
	client := NewClient(&Config{AccessToken: "token"})
	if client.baseURL() != Server {
		t.Errorf("expected %s, got %s", Server, client.baseURL())
	}
}

func TestBaseURL_Custom(t *testing.T) {
	custom := "https://proxy.example.com/api/v1/"
	client := NewClient(&Config{AccessToken: "token", BaseURL: custom})
	if client.baseURL() != custom {
		t.Errorf("expected %s, got %s", custom, client.baseURL())
	}
}

func TestDoApi_UsesCustomBaseURL(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	client := NewClient(&Config{AccessToken: "token", BaseURL: srv.URL + "/api/v1/"})
	client.doApi(context.Background(), http.MethodPost, "sync", url.Values{"sync_token": {"*"}}, nil)

	if gotPath != "/api/v1/sync" {
		t.Errorf("expected request to /api/v1/sync, got %s", gotPath)
	}
}

func TestDoRestApi_UsesCustomBaseURL(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	client := NewClient(&Config{AccessToken: "token", BaseURL: srv.URL + "/api/v1/"})
	client.doRestApi(context.Background(), http.MethodPost, "tasks/quick", nil, nil)

	if gotPath != "/api/v1/tasks/quick" {
		t.Errorf("expected request to /api/v1/tasks/quick, got %s", gotPath)
	}
}
