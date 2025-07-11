package spreadsheet

import (
	"context"
	"fmt"
	"log"
	transaction_domain "money-tracker-bot/internal/domain/transactions"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

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

func (s SpreadsheetService) AppendRow(ctx context.Context, spreadsheetId string, trx transaction_domain.Transaction) {
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
