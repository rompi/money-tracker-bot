package main

import (
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
