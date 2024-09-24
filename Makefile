build:
	go build -o ./bin/geth-modules-sratch

run: build
	./bin/geth-modules-sratch

test: 
	go test -v ./...