.PHONY: all build build-static clean

all: build

build:
	go build -o bin/control_speaker ./cmd/control_speaker
	go build -o bin/mufloctl ./cmd/mufloctl

build-static:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/control_speaker ./cmd/control_speaker
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/mufloctl ./cmd/mufloctl

clean:
	rm -rf bin/
