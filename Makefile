# compile binary

build:
	go build -o bin/tetragon-event-listener pkg/main.go

build-with-version:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/tetragon-event-listener pkg/main.go