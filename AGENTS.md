# Repository Guidelines

## Project Structure & Module Organization
This Go module targets Go 1.23 and follows a hexagonal architecture. The executable entry point lives in `cmd/telebot`, while shared logic sits in `internal/` with adapters (Telegram, Gemini, Google Sheets), domain models, services, and ports kept isolated. Reference material is in `docs/`, and configuration templates live in `.env.example` and `google-service-account.json`.

## Build, Test, and Development Commands
- `make run` (or `go run ./cmd/telebot`) starts the Telegram bot; export `TELEGRAM_BOT_TOKEN`, `GEMINI_API_KEY`, and Sheets credentials first.
- `make build` emits a Linux AMD64 binary at `./bot`.
- `make fmt` and `make lint` wrap `go fmt ./...` and `go vet ./...`.
- `make test` runs `go test -cover ./...`.

## Coding Guidelines (Go)
Format code with `gofmt` (or `make fmt`) and organise imports consistently. Use tabs for indentation, `MixedCaps` for exports, and return wrapped errors with `%w`. Align package layout with the hexagonal layers: domain contracts in `internal/domain`, ports as interfaces under `internal/port`, and adapters implementing those contracts. Prefer dependency injection via constructors and avoid reading environment variables outside `cmd/`.

## Architecture Overview
Hexagonal boundaries must stay intact: the domain layer owns business rules and stays free of external SDKs; ports declare what the domain expects; adapters wrap Telegram, Gemini, and Sheets clients. New integrations should add a port interface, place the adapter under `internal/adapters/<system>`, and wire dependencies in `cmd/telebot/main.go`. Keep DTO translations within adapters to shield the domain from API drift.

## Testing Guidelines
Unit and integration tests reside alongside code as `*_test.go`. Mirror the structure of the code under test (e.g., tests for `internal/service/transactions` sit in the same folder). New features need `go test ./...` to pass without reducing coverage. Prefer table-driven tests and concise fixtures.

## Commit & Pull Request Guidelines
Current history favours prefixing subjects with the feature area, such as `Feature/shopping quota` or `Fix/...`, optionally referencing issues (`(#ID)`). Keep commits focused and written in imperative mood. Pull requests should summarise behaviour changes, cite relevant issues, include configuration or migration notes, and attach test evidence (`make test`) or bot run logs when altering runtime behaviour.

## Security & Configuration Tips
Store all secrets in a local `.env`; do not commit filled `.env` or service account files. When testing Google integrations, prefer stubbed clients from `internal/adapters/google` and redact spreadsheet IDs in logs. Rotate tokens regularly before sharing builds.
