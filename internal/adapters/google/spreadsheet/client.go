package spreadsheet

// Package spreadsheet provides functionality to interact with Google Sheets API for managing transaction data.

import (
	"context"
	"fmt"
	"log"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type CategorySummary struct {
   Category      string
   MonthlyExpenses string
   MonthlyBudget   string
   BudgetLeft      string
}

type SpreadsheetService struct {
	Sheet *sheets.Service
}

func NewSpreadsheetService() *SpreadsheetService {
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile("google-service-account.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return &SpreadsheetService{
		Sheet: srv,
	}

}

func (s SpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) CategorySummary {
	// Add createdAt as UTC+7 timestamp (column G)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("Unable to load Asia/Bangkok timezone: %v", err)
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
		log.Fatalf("Unable to insert data to sheet: %v", err)
	}

	// Fetch summary data from summary sheet
	summaryRange := "summary!A2:D12"
	summaryValues, err := s.Sheet.Spreadsheets.Values.Get(spreadsheetId, summaryRange).Do()
	if err != nil {
		log.Printf("Unable to get monthly budget from summary sheet: %v", err)
		return CategorySummary{}
	}

	// Find the summary for the transaction's category
	var result CategorySummary
	for _, row := range summaryValues.Values {
		if len(row) >= 4 && fmt.Sprintf("%v", row[0]) == trx.Category {
		   result = CategorySummary{
			   Category:         fmt.Sprintf("%v", row[0]),
			   MonthlyExpenses:  fmt.Sprintf("%v", row[1]),
			   MonthlyBudget:    fmt.Sprintf("%v", row[2]),
			   BudgetLeft:       fmt.Sprintf("%v", row[3]),
		   }
			break
		}
	}
	// Optionally handle missing category
	return result
}

func (s SpreadsheetService) GetCellValue(ctx context.Context, spreadsheetId string) {
	values, err := s.Sheet.Spreadsheets.Values.Get(spreadsheetId, "Sheet1!A2:E7").Do()

	if err != nil {
		log.Fatalf("Unable to Get data from sheet: %v", err)
	}

	for _, value := range values.Values {
		fmt.Println(value)
	}
}
