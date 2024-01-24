GO := go
BIN := echoimages

.PHONY: all clean

all: build

build:
	@$(GO) build -o $(BIN) $(PKG)

run: build
	@ENV="dev" ./$(BIN)

clean:
	@$(GO) clean $(PKG)
