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

func TestMain_NoArguments_PrintsUsageAndExitsNonZero(t *testing.T) {
	cmd := exec.Command(binPath)
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected a non-zero exit status when no arguments are given")
	}
	if !strings.Contains(string(out), "Usage:") {
		t.Errorf("output = %q, want it to contain %q", string(out), "Usage:")
	}
}