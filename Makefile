.PHONY: run

debug:
	go run cmd/main.go debug

run:
	go run cmd/main.go apply

all: run
