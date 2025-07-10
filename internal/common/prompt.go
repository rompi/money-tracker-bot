package common

import (
	"fmt"
	"strings"
)

// TransactionCategoryList is the static list of allowed categories for transactions
var TransactionCategoryList = []string{
	"Groceries",
	"Utilities",
	"Entertainment",
	"Gifting",
	"Household",
	"Eating Out",
	"Health",
	"Transportation",
	"Savings",
	"Emergency",
	"Rent House",
}

// SourceAccountList is the static list of allowed source accounts
var SourceAccountList = []string{
	"GOPAY",
	"BCA",
	"OVO",
	"DANA",
	"ISAKU",
	"MANDIRI",
	"BNI",
	"BRI",
	"CASH",
}

// PromptParams holds parameters for building the prompt
// If IsImage is true, FileID must be set. If false, Message and CurrentDate must be set.
type PromptParams struct {
	IsImage     bool
	FileID      string
	Message     string
	CurrentDate string
}

// BuildPrompt builds the prompt for Gemini based on the input params
func BuildPrompt(params PromptParams) string {
	categoryStr := strings.Join(TransactionCategoryList, " / ")
	sourceAccountStr := strings.Join(SourceAccountList, " / ")

	fields := fmt.Sprintf(`Fields:
  - title (summary of the transaction notes)
  - transaction_date (format always YYYY-MM-DD)
  - amount (in rupiah, format 1,000,000 for transaction with amont 1 million. if it is 100k then output should be 100,000)
  - notes (details of the transaction, containing items bought)
  - category (%s)`,
		categoryStr)

	if params.IsImage {
		fields += fmt.Sprintf(`
  - destination_number
  - source_account (only %s)
  - file_id %s`, sourceAccountStr, params.FileID)
	} else {
		fields += `
  - file_id should be empty`
	}

	var inputDesc string
	if params.IsImage {
		inputDesc = "from the image"
	} else {
		inputDesc = fmt.Sprintf("from the following message: %s", params.Message)
	}

	var dateLine string
	if params.IsImage {
		dateLine = ""
	} else {
		dateLine = fmt.Sprintf("  - transaction_date should be %s (format always YYYY-MM-DD)\n", params.CurrentDate)
	}

	exampleFileID := params.FileID
	if !params.IsImage {
		exampleFileID = ""
	}

	prompt := fmt.Sprintf(`Please extract the following data %s and return it as valid JSON.

%s
%sIMPORTANT:
Respond ONLY with raw JSON.
No explanation, no formatting, no code blocks.

Example:
{
  "title": "Transfer to ABC Cafe",
  "transaction_date": "2025-03-30",
  "amount": "150",
  "notes": "Lunch at ABC cafe",
  "destination_number": "0524012911",
  "source_account": "Gopay",
  "category": "Groceries",
  "file_id": "%s"
}`,
		inputDesc, fields, dateLine, exampleFileID)
	return prompt
}
