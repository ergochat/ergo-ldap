.PHONY: all build test gofmt

all: build

build:
	go build -v .

test:
	cd ldap && go test . && go vet .
	go vet .
	./.check-gofmt.sh

gofmt:
	./.check-gofmt.sh --fix
