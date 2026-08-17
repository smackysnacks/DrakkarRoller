default:
    just --list

# Build the roller binary
build:
    go build -o roller .

# Run the roller against a local game instance
run:
    go run .

# Format all Go source files
fmt:
    gofmt -l -w .

# Run go vet
vet:
    go vet ./...

# Tidy go.mod/go.sum
tidy:
    go mod tidy

# Run fmt, vet and tests
check: fmt vet
    go test ./...

# Remove build artifacts and logs
clean:
    rm -f roller roller.exe log.txt
