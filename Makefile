# Makefile
APP_NAME=leave
CMD_PATH=./

.PHONY: run build clean

run:
	go run $(CMD_PATH)

build:
	go build -o $(APP_NAME) $(CMD_PATH)

dev: build
	./$(APP_NAME)

clean:
	rm -f $(APP_NAME)

install:
	go install $(CMD_PATH)