package spreadsheet

type MockSpreadsheetService struct{}

func (m *MockSpreadsheetService) AppendRow(ctx interface{}, spreadsheetId string, trx interface{}) {}
func (m *MockSpreadsheetService) GetCellValue(ctx interface{}, spreadsheetId string)               {}

// Add more mock methods as needed for your tests
