build:
	@go build -o bin/go-pay ./main.go

test:
	@go test -v ./...
	
run: build
	@./bin/go-pay