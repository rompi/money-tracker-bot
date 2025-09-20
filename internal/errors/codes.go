package errors

// Error codes for categorizing different types of application errors
const (
	// Configuration and startup errors
	ErrCodeConfig = "CONFIG_ERROR"

	// External service errors
	ErrCodeTelegram    = "TELEGRAM_ERROR"
	ErrCodeGemini      = "GEMINI_ERROR"
	ErrCodeSpreadsheet = "SPREADSHEET_ERROR"

	// Internal operation errors
	ErrCodeFileOperation = "FILE_ERROR"
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeTransaction   = "TRANSACTION_ERROR"

	// Network and connectivity errors
	ErrCodeNetwork = "NETWORK_ERROR"
	ErrCodeTimeout = "TIMEOUT_ERROR"

	// Data and persistence errors
	ErrCodeDataAccess    = "DATA_ACCESS_ERROR"
	ErrCodeDataFormat    = "DATA_FORMAT_ERROR"
	ErrCodeDataIntegrity = "DATA_INTEGRITY_ERROR"
)

// Error severity levels
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}
