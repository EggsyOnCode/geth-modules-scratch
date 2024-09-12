build:
	go build -o ./bin/go-audio-codec

run: build
	./bin/go-audio-codec

test: 
	go test -v ./...