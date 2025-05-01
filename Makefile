build: 
    @go build -o gobank
run: build
    @./bin/gobank
test:
    @go test -v  ./...