APPNAME=$(shell basename $(CURDIR))

default:
	go build -o bin/$(APPNAME)

fmt:
	gofmt -w .

.PHONY: default
