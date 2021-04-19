SHELL=/usr/bin/env bash

.PHONY: clean
clean:
	rm beacon-plus

.PHONY: all
all:
	go build -o beacon-plus *.go

.PHONY: calibnet
calibnet:
	go build -tags=calibnet -o beacon-plus *.go