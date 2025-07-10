package common

import (
	"strings"
	"testing"
)

func TestBuildPrompt_Image(t *testing.T) {
	params := PromptParams{
		IsImage: true,
		FileID:  "testfile.jpg",
	}
	prompt := BuildPrompt(params)

	if !strings.Contains(prompt, "from the image") {
		t.Errorf("Prompt should mention 'from the image'")
	}
	if !strings.Contains(prompt, "testfile.jpg") {
		t.Errorf("Prompt should include the file ID")
	}
	if !strings.Contains(prompt, "source_account (only ") {
		t.Errorf("Prompt should include source_account list")
	}
	if !strings.Contains(prompt, "category (") {
		t.Errorf("Prompt should include category list")
	}
}

func TestBuildPrompt_Text(t *testing.T) {
	params := PromptParams{
		IsImage:     false,
		Message:     "Transfer 100k to Budi",
		CurrentDate: "2025-07-10",
	}
	prompt := BuildPrompt(params)

	if !strings.Contains(prompt, "from the following message: Transfer 100k to Budi") {
		t.Errorf("Prompt should mention the message")
	}
	if !strings.Contains(prompt, "- transaction_date should be 2025-07-10") {
		t.Errorf("Prompt should include the current date")
	}
	if !strings.Contains(prompt, "file_id should be empty") {
		t.Errorf("Prompt should specify file_id should be empty")
	}
	if !strings.Contains(prompt, "category (") {
		t.Errorf("Prompt should include category list")
	}
}
