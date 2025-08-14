package main

import (
	"os"
	"testing"
)

func TestStartBotWithDeps_MissingEnvVars(t *testing.T) {
	err := startBotWithDeps("", "", nil, nil)
	if err == nil {
		t.Error("expected error when env vars are missing, got nil")
	}
}

func TestStartBotWithDeps_AllEnvVarsPresent(t *testing.T) {
	telegramToken := "dummy-token"
	apiKey := "dummy-key"
	// Use empty struct for mocks, as interfaces are now empty
	mockSpreadsheet := struct{}{}
	mockGemini := struct{}{}
	err := startBotWithDeps(telegramToken, apiKey, mockSpreadsheet, mockGemini)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestErrEnvVarMissing_Error(t *testing.T) {
	err := ErrEnvVarMissing("TELEGRAM_BOT_TOKEN")
	got := err.Error()
	want := "required environment variable not set: TELEGRAM_BOT_TOKEN"
	if got != want {
		t.Errorf("expected '%s', got '%s'", want, got)
	}
}

func TestStartBot_MissingEnvVars(t *testing.T) {
	// Save and restore env
	origToken := getenv("TELEGRAM_BOT_TOKEN")
	origKey := getenv("GEMINI_API_KEY")
	os.Setenv("TELEGRAM_BOT_TOKEN", "")
	os.Setenv("GEMINI_API_KEY", "")
	defer func() {
		os.Setenv("TELEGRAM_BOT_TOKEN", origToken)
		os.Setenv("GEMINI_API_KEY", origKey)
	}()
	testBotDeps.Override = true
	testBotDeps.SpreadsheetService = struct{}{}
	testBotDeps.GeminiClient = struct{}{}
	err := startBot()
	testBotDeps.Override = false
	if err == nil {
		t.Error("expected error when env vars are missing, got nil")
	}
}

func TestStartBot_AllEnvVarsPresent(t *testing.T) {
	origToken := getenv("TELEGRAM_BOT_TOKEN")
	origKey := getenv("GEMINI_API_KEY")
	os.Setenv("TELEGRAM_BOT_TOKEN", "dummy-token")
	os.Setenv("GEMINI_API_KEY", "dummy-key")
	defer func() {
		os.Setenv("TELEGRAM_BOT_TOKEN", origToken)
		os.Setenv("GEMINI_API_KEY", origKey)
	}()
	testBotDeps.Override = true
	testBotDeps.SpreadsheetService = struct{}{}
	testBotDeps.GeminiClient = struct{}{}
	err := startBot()
	testBotDeps.Override = false
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

// getenv is a helper to avoid panic if env is not set
func getenv(key string) string {
	v := os.Getenv(key)
	return v
}
