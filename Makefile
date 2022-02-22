.PHONY: run
run:
	go run cmd/gophermart/main.go

.PHONY: test
test:
	go test -v -count=1 ./...