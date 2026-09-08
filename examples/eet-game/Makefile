.PHONY: all build run test tidy clean

BINARY := bin/eet-game
CMD     := ./cmd/server

all: build

build:
	@mkdir -p bin
	go build -o $(BINARY) $(CMD)
	@echo "✅ Built $(BINARY)"

run:
	go run $(CMD)

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/
	@echo "🧹 Cleaned"
