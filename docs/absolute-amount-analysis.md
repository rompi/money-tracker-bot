# Absolute Amount Analysis

## Background
Currently, when processing transaction amounts from user input, the system may receive negative or positive numbers. We want to ensure that regardless of the input format, the system always processes the amount as an absolute (positive) number.

## Current Behavior
- Users can input amounts with or without a negative sign
- The system processes these amounts as-is, maintaining the sign
- The amount handling is primarily done in two places:
  1. Initial prompt generation (`internal/common/prompt.go`)
  2. Response parsing (`internal/adapters/gemini/gemini.go`)

## Requirement
- All transaction amounts should be processed as absolute (positive) numbers
- This should happen during the initial prompt/input processing stage
- The sign of the transaction (whether it's expense or income) should be determined by the context/category, not by the numerical input

## Implementation Areas to Modify

1. **Prompt Processing**
   - Location: `internal/common/prompt.go`
   - Update prompt template to explicitly instruct AI to return absolute values
   - Modify example in the prompt to demonstrate absolute value usage
   - Add clear instructions about:
     - Always returning positive amounts
     - Using context (not number sign) to determine transaction type

2. **Gemini Adapter**
   - Location: `internal/adapters/gemini/gemini.go`
   - Add preprocessing for transaction amounts in both methods:
     - `ReadImageToTransaction`
     - `TextToTransaction`
   - Add validation to ensure amounts are always positive
   - Implement fallback conversion to absolute value if needed

## Test Cases to Consider
1. Input with positive number:
   - Text: "spent 100 on groceries"
   - Expected: amount = "100"
2. Input with negative number:
   - Text: "spent -100 on groceries"
   - Expected: amount = "100"
3. Input with currency symbols:
   - Text: "spent $100 on groceries"
   - Expected: amount = "100"
4. Input with decimal points:
   - Text: "spent 100.50 on groceries"
   - Expected: amount = "100.50"
5. Input with transaction context:
   - Text: "earned 100 from freelancing"
   - Expected: amount = "100" (positive, income determined by context)

## Implementation Steps
1. Modify `internal/common/prompt.go`:
   ```
   - Update field description for amount to specify absolute values
   - Add explicit instruction about using context for transaction type
   - Update example JSON to demonstrate proper format
   ```

2. Update `internal/adapters/gemini/gemini.go`:
   ```
   - Add amount validation after JSON parsing
   - Ensure absolute value conversion if needed
   - Apply changes to both image and text processing
   ```

3. Add test cases:
   ```
   - Test positive input amounts
   - Test negative input amounts
   - Test amounts with currency symbols
   - Test decimal amounts
   - Test different transaction contexts (spending vs earning)
   ```

## Success Criteria
1. Input Processing:
   - All amounts are stored as positive numbers
   - Currency symbols are properly handled
   - Decimal values are preserved
   - Negative signs are removed

2. Context Handling:
   - Transaction type is determined by context/category
   - Words like "spent", "bought", "earned", "received" are used to determine type
   - Sign of the number doesn't affect categorization

3. System Integrity:
   - Existing functionality remains intact
   - All test cases pass
   - No regression in current features
