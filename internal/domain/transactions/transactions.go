package transaction_domain

type Transaction struct {
   TransactionDate   string `json:"transaction_date"`
   Amount            string `json:"amount"`
   AmountCurrency    string `json:"amount_currency"`
   Notes             string `json:"notes"`
   DestinationName   string `json:"destination_name"`
   DestinationNumber string `json:"destination_number"`
   SourceAccount     string `json:"source_account"`
   Category          string `json:"category"`
   Title             string `json:"title"`
   FileID            string `json:"file_id"`
   CreatedBy         string `json:"created_by"`
   WarningMessage    string `json:"warning_message,omitempty"`
}
