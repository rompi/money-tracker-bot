# Gemini AI Adapter

## Package: `internal/adapters/gemini`

### Purpose
Adapter for Google Gemini AI service that converts images and text into structured transaction data using generative AI.

### Key Components

#### `gemini.go`
- **Purpose**: Gemini AI client for transaction data extraction
- **Key Structures**:
  - `GeminiClient`: Main client with Gemini API integration
  - `GenerativeModelPort`: Interface for testability
- **Key Functions**:
  - `ReadImageToTransaction()`: Processes receipt/transaction images into structured data
  - `TextToTransaction()`: Converts text messages into transaction records
  - `GenerateContent()`: Low-level Gemini API interaction

#### AI Processing Flow
1. **Image Processing**:
   - Reads image files (JPEG format)
   - Sends to Gemini with structured prompt
   - Extracts transaction details (amount, category, notes, etc.)
   - Ensures positive amounts
   - Cleans up temporary files

2. **Text Processing**:
   - Processes natural language transaction descriptions
   - Uses current date as transaction date
   - Extracts structured data using AI prompts

#### Data Validation
- **Amount Normalization**: Ensures all amounts are positive
- **JSON Parsing**: Handles AI response cleanup and parsing
- **Error Handling**: Graceful handling of AI response failures

#### Dependencies
- Google Generative AI Go SDK (`github.com/google/generative-ai-go/genai`)
- Common prompt building utilities
- Transaction domain models

### AI Model Configuration
- Uses `gemini-2.0-flash` model for optimal performance
- Structured prompts with predefined categories and accounts
- JSON-only responses for reliable parsing