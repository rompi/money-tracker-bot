package spreadsheet

import (
	"context"
	"fmt"
	"log"
	transaction_domain "rompi/gobot/internal/domain/transactions"

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

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			trx.TransactionDatetime,
			trx.Title,
			trx.FileID,
			trx.Amount,
			trx.Category,
			trx.Notes,
			trx.DestinationName,
			trx.SourceAccount,
			trx.CreatedBy,
		}},
	}

	_, err := s.Sheet.Spreadsheets.Values.Append(spreadsheetId, "Expense Tracker!A:I", values).ValueInputOption("USER_ENTERED").Do()

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

func (s SpreadsheetService) Update(ctx context.Context, spreadsheetId string) {

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			"Japan",
			"Software Engineer Lead",
		}},
	}

	_, err := s.Sheet.Spreadsheets.Values.Update(spreadsheetId, "Sheet1!B2:C2", values).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		log.Fatalf("Unable to insert data to sheet: %v", err)
	}

}
