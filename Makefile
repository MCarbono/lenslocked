default: run

run:
	go run main.go

tests:
	go test ./... -v

gen:
	go generate ./...