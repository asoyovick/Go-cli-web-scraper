package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var binPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "go-cli-web-scraper-bin")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	binPath = filepath.Join(dir, "go-cli-web-scraper")
	if runtimeIsWindows() {
		binPath += ".exe"
	}

	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build CLI binary: " + err.Error() + "\n" + string(out))
	}

	os.Exit(m.Run())
}

func runtimeIsWindows() bool {
	return os.PathSeparator == '\\'
}

func TestMain_NoArgs(t *testing.T) {
	cmd := exec.Command(binPath)
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected a non-zero exit status when no arguments are given")
	}
	if !strings.Contains(string(out), "Usage:") {
		t.Errorf("output = %q, want it to contain %q", string(out), "Usage:")
	}
}
func TestMain_LocalFile(t *testing.T) {
	dir := t.TempDir()
	fixture := filepath.Join(dir, "page.html")
	content := `<html><head><title>CLI Fixture</title></head><body>
		<a href="/one">One</a>
		<a href="/two">Two</a>
		<img src="/pic.png">
	</body></html>`
	if err := os.WriteFile(fixture, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	cmd := exec.Command(binPath, fixture)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error running CLI: %v\noutput: %s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "Title: CLI Fixture") {
		t.Errorf("output missing title line, got:\n%s", output)
	}
	if !strings.Contains(output, "Links Found(2):") {
		t.Errorf("output missing expected link count, got:\n%s", output)
	}
	if !strings.Contains(output, "Images Found (1):") {
		t.Errorf("output missing expected image count, got:\n%s", output)
	}
}

func TestMain_RemoteURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Remote Fixture</title></head><body>
			<a href="/a">a</a>
		</body></html>`))
	}))
	defer srv.Close()

	cmd := exec.Command(binPath, srv.URL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error running CLI: %v\noutput: %s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "Title: Remote Fixture") {
		t.Errorf("output missing title line, got:\n%s", output)
	}
	if !strings.Contains(output, srv.URL+"/a") {
		t.Errorf("output missing resolved link, got:\n%s", output)
	}
}
func TestMain_MissingFile(t *testing.T) {
	cmd := exec.Command(binPath, "/definitely/does/not/exist.html")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit code 0 even on scrape error, got err: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), "Error:") {
		t.Errorf("output = %q, want it to contain %q", string(out), "Error:")
	}
}

func TestMain_BadURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	badURL := srv.URL
	srv.Close()

	cmd := exec.Command(binPath, badURL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit code 0 even on network error, got err: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), "Error:") {
		t.Errorf("output = %q, want it to contain %q", string(out), "Error:")
	}
}

func TestMain_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // comfortably longer than the client's -timeout below
	}))
	defer srv.Close()

	cmd := exec.Command(binPath, "-timeout=50ms", srv.URL)
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "Error:") {
		t.Errorf("output = %q, want a timeout error", string(out))
	}
}

func TestMain_InvalidTimeout(t *testing.T) {
	cmd := exec.Command(binPath, "-timeout=notaduration", "https://example.com")
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected a non-zero exit status for an invalid -timeout value")
	}
}

func TestMain_EmptyPage(t *testing.T) {
	dir := t.TempDir()
	fixture := filepath.Join(dir, "empty.html")
	if err := os.WriteFile(fixture, []byte(`<html><head><title>Empty</title></head><body></body></html>`), 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	cmd := exec.Command(binPath, fixture)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "Links Found(0):") {
		t.Errorf("output missing zero link count, got:\n%s", output)
	}
	if !strings.Contains(output, "Images Found (0):") {
		t.Errorf("output missing zero image count, got:\n%s", output)
	}
}