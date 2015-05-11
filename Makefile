APIPORT ?= 8888

.PHONY: run
run:
	APIPORT=$(APIPORT) \
	go run cmd/wiretap/main.go
