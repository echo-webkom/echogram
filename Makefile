GO := go
BIN := echoimages

.PHONY: all clean

all: build

build:
	@$(GO) build -o $(BIN) $(PKG)

run: build
	@./$(BIN)

clean:
	@$(GO) clean $(PKG)
