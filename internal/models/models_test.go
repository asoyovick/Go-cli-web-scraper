package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPageData_JSON(t *testing.T) {
	pd := PageData{
		URL:        "https://example.com",
		Title:      "Example",
		Links:      []string{"https://example.com/a", "https://example.com/b"},
		Images:     []string{"https://example.com/img.png"},
		StatusCode: 200,
	}

	data, err := json.Marshal(pd)
	if err != nil {
		t.Fatalf("unexpected error marshaling PageData: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unexpected error unmarshaling result: %v", err)
	}

	wantKeys := []string{"url", "title", "links", "images", "status_code"}
	for _, k := range wantKeys {
		if _, ok := out[k]; !ok {
			t.Errorf("expected JSON key %q to be present, got keys %v", k, out)
		}
	}

	if out["url"] != "https://example.com" {
		t.Errorf("url = %v, want %q", out["url"], "https://example.com")
	}
	if out["title"] != "Example" {
		t.Errorf("title = %v, want %q", out["title"], "Example")
	}
	if out["status_code"].(float64) != 200 {
		t.Errorf("status_code = %v, want 200", out["status_code"])
	}
}

func TestPageData_Zero(t *testing.T) {
	var pd PageData
	data, err := json.Marshal(pd)
	if err != nil {
		t.Fatalf("unexpected error marshaling zero-value PageData: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unexpected error unmarshaling result: %v", err)
	}

	if out["url"] != "" {
		t.Errorf("url = %v, want empty string", out["url"])
	}
	if out["links"] != nil {
		t.Errorf("links = %v, want nil/null", out["links"])
	}
	if out["status_code"].(float64) != 0 {
		t.Errorf("status_code = %v, want 0", out["status_code"])
	}
}

func TestPageData_RoundTrip(t *testing.T) {
	original := PageData{
		URL:        "file:///tmp/page.html",
		Title:      "Round Trip",
		Links:      []string{"a", "b", "c"},
		Images:     nil,
		StatusCode: 404,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded PageData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.URL != original.URL {
		t.Errorf("URL = %q, want %q", decoded.URL, original.URL)
	}
	if decoded.Title != original.Title {
		t.Errorf("Title = %q, want %q", decoded.Title, original.Title)
	}
	if len(decoded.Links) != len(original.Links) {
		t.Errorf("Links = %v, want %v", decoded.Links, original.Links)
	}
	if decoded.StatusCode != original.StatusCode {
		t.Errorf("StatusCode = %d, want %d", decoded.StatusCode, original.StatusCode)
	}
}