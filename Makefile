# Parent Makefiles https://github.com/c4s4/make

include ~/.make/golang.mk

GOOSARCH = $(shell go tool dist list | grep -v android | grep -v plan9)

test:    go-test    # Run unit tests
release: go-release # Perform release
