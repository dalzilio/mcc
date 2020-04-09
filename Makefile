
VERSION := $(shell git describe --tags --dirty --always)

DATE := $(shell date -u +%d/%m/%Y)

main:
	go install -ldflags="-X 'github.com/dalzilio/mcc/cmd.version=$(VERSION)' -X 'github.com/dalzilio/mcc/cmd.builddate=$(DATE)'"

build:
	go build -ldflags="-X 'github.com/dalzilio/mcc/cmd.version=$(VERSION)' -X 'github.com/dalzilio/mcc/cmd.builddate=$(DATE)'"
