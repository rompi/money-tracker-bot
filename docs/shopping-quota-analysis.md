# Shopping Quota Feature Analysis

## Objective
Add new properties `quota` and `quota left` to the telegram response when a user records a new transaction through the bot.

## Cur## Testing Requirements and Scenarios

### 1. Unit Tests

#### a) CategorySummary Data Structure (`client_test.go`)
```go
func TestCategorySummary_WithQuota(t *testing.T) {
    cases := []struct {
        name     string
        input    []interface{}
        expected CategorySummary
    }{
        {
            name:  "Complete data with quota",
            input: []interface{}{"Food", "1000", "5000", "4000", "2000", "1500"},
            expected: CategorySummary{
                Category:        "Food",
                MonthlyExpenses: "1000",
                MonthlyBudget:   "5000",
                BudgetLeft:      "4000",
                Quota:           "2000",
                QuotaLeft:       "1500",
            },
        },
        {
            name:  "Missing quota data",
            input: []interface{}{"Food", "1000", "5000", "4000"},
            expected: CategorySummary{
                Category:        "Food",
                MonthlyExpenses: "1000",
                MonthlyBudget:   "5000",
                BudgetLeft:      "4000",
                Quota:           "",
                QuotaLeft:       "",
            },
        },
    }
    // Test implementation
}
```

#### b) Message Formatting (`handler_test.go`)
```go
func TestMessageFormatting_WithQuota(t *testing.T) {
    cases := []struct {
        name     string
        summary  CategorySummary
        expected string
    }{
        {
            name: "Complete quota information",
            summary: CategorySummary{
                Quota:     "2000",
                QuotaLeft: "1500",
            },
            expected: "Monthly Quota: 2000\nQuota Left: 1500",
        },
        {
            name: "Missing quota information",
            summary: CategorySummary{
                Quota:     "",
                QuotaLeft: "",
            },
            expected: "Monthly Quota: -\nQuota Left: -",
        },
    }
    // Test implementation
}
```

### 2. Integration Tests

a) End-to-end flow testing:
- Record transaction
- Verify spreadsheet read range
- Validate response message format
- Check quota data accuracy

b) Edge case testing:
- Categories without quota
- Zero or negative quota amounts
- Categories not found in summary sheet
- Missing spreadsheet columns

### 3. Manual Testing Scenarios

a) Transaction recording:
1. Record transaction in category with quota
2. Verify response shows correct quota information
3. Check quota left calculation accuracy

b) Edge cases:
1. Record transaction in category without quota
2. Record transaction exceeding quota
3. Record multiple transactions affecting quota

## Technical Considerationsmentation Analysis
Currently, when a transaction is recorded:
1. The transaction data is stored in the `detailed` sheet
2. The system fetches summary data from `summary!A2:D12` which includes:
   - Category (Column A)
   - Monthly Expenses (Column B)
   - Monthly Budget (Column C)
   - Budget Left (Column D)

## Required Changes by File

### 1. `/internal/adapters/google/spreadsheet/client.go`

a) Update `CategorySummary` struct:
```go
type CategorySummary struct {
    Category        string
    MonthlyExpenses string
    MonthlyBudget   string
    BudgetLeft      string
    Quota           string  // New field
    QuotaLeft       string  // New field
}
```

b) Modify `AppendRow` method:
- Update summary range from `summary!A2:D12` to `summary!A2:F12`
- Update result mapping to include new fields:
```go
result = CategorySummary{
    Category:        fmt.Sprintf("%v", row[0]),
    MonthlyExpenses: fmt.Sprintf("%v", row[1]),
    MonthlyBudget:   fmt.Sprintf("%v", row[2]),
    BudgetLeft:      fmt.Sprintf("%v", row[3]),
    Quota:           fmt.Sprintf("%v", row[4]),  // New
    QuotaLeft:       fmt.Sprintf("%v", row[5]),  // New
}
```

### 2. `/internal/adapters/telegram/handler.go`

a) Update message formatting in both `handlePhoto` and `handleMessage` methods:
```go
msgText := fmt.Sprintf(
    "Saved %s ✅\nCategory: %s\nAmount: %s\nNotes: %s\nLink: %s\n"+
    "Monthly Expenses: %s\nMonthly Budget: %s\nBudget Left: %s\n"+
    "Monthly Quota: %s\nQuota Left: %s",  // New lines
    fileType, // "photo" or "text"
    transaction.Category,
    rupiah,
    transaction.Notes,
    spreadsheetLink,
    summary.MonthlyExpenses,
    summary.MonthlyBudget,
    summary.BudgetLeft,
    summary.Quota,        // New
    summary.QuotaLeft,    // New
)
```

### 3. Test Files to Update

a) `/internal/adapters/google/spreadsheet/client_test.go`:
- Add test cases for quota fields in CategorySummary
- Update mock data to include quota columns
- Test edge cases for missing quota data

b) `/internal/adapters/telegram/handler_test.go`:
- Update mock responses to include quota information
- Add test cases for quota display in messages
- Test formatting with and without quota data

### 4. Mock Files to Update

a) `/internal/adapters/google/spreadsheet/mock.go`:
- Update mock implementation to include quota fields
- Add test scenarios for quota data handling

## Technical Implementation Details

### 1. Spreadsheet Service Changes
- Update the range in `AppendRow` method to include columns E and F from the summary sheet
- Current range: `summary!A2:D12`
- New range: `summary!A2:F12`
- Modify `CategorySummary` struct to include new fields:
  ```go
  type CategorySummary struct {
      Category        string
      MonthlyExpenses string
      MonthlyBudget   string
      BudgetLeft      string
      Quota           string  // New field
      QuotaLeft       string  // New field
  }
  ```

### 2. Transaction Service Changes
The transaction service will need to include the new quota information in the response message format.

### 3. Message Format Changes
Current format includes category summary with expenses and budget. New format will add quota information:
```
Transaction recorded:
Amount: [amount]
Category: [category]
Notes: [notes]

Category Summary:
Monthly Expenses: [expenses]
Monthly Budget: [budget]
Budget Left: [budget_left]
Monthly Quota: [quota]       <- New
Quota Left: [quota_left]    <- New
```

## Implementation Steps

### 1. Update Data Structure
a) In `/internal/adapters/google/spreadsheet/client.go`:
```go
// Update CategorySummary struct to include quota fields
type CategorySummary struct {
    Category        string
    MonthlyExpenses string
    MonthlyBudget   string
    BudgetLeft      string
    Quota           string
    QuotaLeft       string
}
```

b) Update mock implementations in test files to include new fields.

### 2. Update Spreadsheet Integration
a) In `AppendRow` method:
```go
// Update range to include quota columns
summaryRange := "summary!A2:F12"
```

b) Update result mapping to include quota fields:
```go
result = CategorySummary{
    Category:        fmt.Sprintf("%v", row[0]),
    MonthlyExpenses: fmt.Sprintf("%v", row[1]),
    MonthlyBudget:   fmt.Sprintf("%v", row[2]),
    BudgetLeft:      fmt.Sprintf("%v", row[3]),
    Quota:           fmt.Sprintf("%v", row[4]),
    QuotaLeft:       fmt.Sprintf("%v", row[5]),
}
```

### 3. Update Message Formatting
a) In `/internal/adapters/telegram/handler.go`:
- Modify both `handlePhoto` and `handleMessage` methods
- Add quota information to response messages
- Update message formatting to include new fields

### 4. Update Tests
a) Add test cases in `client_test.go`:
- Test quota data retrieval
- Test handling of missing quota data
- Test edge cases

b) Update handler tests:
- Test message formatting with quota data
- Test handling of missing quota information

### 5. Documentation
a) Update technical documentation to include:
- New spreadsheet column requirements
- New response format
- Updated data structure

## Testing Requirements and Scenarios
1. Verify quota data is correctly fetched from summary sheet
2. Ensure response messages include quota information
3. Test with different categories to ensure quota information is category-specific
4. Test edge cases:
   - Categories without quota
   - Zero or negative quota amounts
   - Categories not found in summary sheet

## Technical Considerations
- No database schema changes required
- Changes are isolated to application layer
- Backward compatible with existing spreadsheet structure
- Minimal impact on existing functionality

## Dependencies
- Access to columns E and F in summary sheet must be granted
- Summary sheet must maintain consistent structure for quota data
