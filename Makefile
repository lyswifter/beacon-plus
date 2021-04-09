SHELL=/usr/bin/env bash

.PHONY: clean
clean:
	rm beacon-plus

.PHONY: all
all:
	go build -o beacon-plus *.go
