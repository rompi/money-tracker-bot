package spreadsheet

// Package spreadsheet provides functionality to interact with Google Sheets API for managing transaction data.

import (
	"context"
	"fmt"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	"money-tracker-bot/internal/errors"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type CategorySummary struct {
	Category        string
	MonthlyExpenses string
	MonthlyBudget   string
	BudgetLeft      string
	Quota           string
	QuotaLeft       string
}

type SpreadsheetService struct {
	Sheet *sheets.Service
}

func NewSpreadsheetService() (*SpreadsheetService, error) {
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile("google-service-account.json"))
	if err != nil {
		return nil, errors.NewSpreadsheetCriticalError("failed to create Google Sheets client", err).
			WithContext("credentials_file", "google-service-account.json").
			WithComponent("spreadsheet-client")
	}

	return &SpreadsheetService{
		Sheet: srv,
	}, nil
}

func (s SpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) (CategorySummary, error) {
	// Add createdAt as UTC+7 timestamp (column G)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return CategorySummary{}, errors.NewConfigError("failed to load timezone", err).
			WithContext("timezone", "Asia/Bangkok").
			WithComponent("spreadsheet-client")
	}
	createdAt := time.Now().In(loc).Format("2006-01-02 15:04:05")

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			trx.TransactionDate,
			trx.Category,
			"",
			trx.Notes,
			trx.Amount,
			trx.CreatedBy,
			trx.FileID,
			createdAt,
		}},
	}

	// Update range to include column G
	_, err = s.Sheet.Spreadsheets.Values.Append(spreadsheetId, "detailed!A:G", values).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return CategorySummary{}, errors.NewSpreadsheetError("failed to insert data to sheet", err).
			WithContext("spreadsheet_id", spreadsheetId).
			WithContext("range", "detailed!A:G").
			WithComponent("spreadsheet-client")
	}

	// Fetch summary data from summary sheet (now includes columns E and F)
	summaryRange := "summary!A2:F12"
	summaryValues, err := s.Sheet.Spreadsheets.Values.Get(spreadsheetId, summaryRange).Do()
	if err != nil {
		// Log warning but don't fail the transaction append
		errorWithContext := errors.NewSpreadsheetError("failed to get summary data", err).
			WithContext("spreadsheet_id", spreadsheetId).
			WithContext("range", summaryRange).
			WithComponent("spreadsheet-client")
		errors.HandleError(errorWithContext, "retrieving category summary")
		return CategorySummary{}, nil
	}

	// Find the summary for the transaction's category
	var result CategorySummary
	for _, row := range summaryValues.Values {
		if len(row) >= 4 && fmt.Sprintf("%v", row[0]) == trx.Category {
			// Defensive: handle missing quota columns gracefully
			quota := ""
			quotaLeft := ""
			if len(row) > 4 {
				quota = fmt.Sprintf("%v", row[4])
			}
			if len(row) > 5 {
				quotaLeft = fmt.Sprintf("%v", row[5])
			}
			result = CategorySummary{
				Category:        fmt.Sprintf("%v", row[0]),
				MonthlyExpenses: fmt.Sprintf("%v", row[1]),
				MonthlyBudget:   fmt.Sprintf("%v", row[2]),
				BudgetLeft:      fmt.Sprintf("%v", row[3]),
				Quota:           quota,
				QuotaLeft:       quotaLeft,
			}
			break
		}
	}
	// Optionally handle missing category
	return result, nil
}

func (s SpreadsheetService) GetCellValue(ctx context.Context, spreadsheetId string) error {
	values, err := s.Sheet.Spreadsheets.Values.Get(spreadsheetId, "Sheet1!A2:E7").Do()

	if err != nil {
		return errors.NewSpreadsheetError("failed to get cell values", err).
			WithContext("spreadsheet_id", spreadsheetId).
			WithContext("range", "Sheet1!A2:E7").
			WithComponent("spreadsheet-client")
	}

	for _, value := range values.Values {
		fmt.Println(value)
	}
	return nil
}
