test:
	go vet ./...
	go get github.com/bmizerany/assert
	go test -cover -short ./...
