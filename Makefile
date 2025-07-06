.PHONY: build run clean sqlc-generate db-reset

# Build the application
build:
	go build -o goautomation .

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -f goautomation
	rm -rf data/

# Generate Go code from SQL queries
sqlc-generate:
	sqlc generate

# Reset database (delete and recreate)
db-reset:
	rm -rf data/
	mkdir -p data

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Show database schema
db-schema:
	sqlite3 data/automation.sqlite ".schema"

# Show recent logs from database
db-logs:
	sqlite3 data/automation.sqlite "SELECT * FROM automation_logs ORDER BY timestamp DESC LIMIT 10;"
