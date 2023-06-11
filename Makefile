default: run

run:
	go run main.go

test:
	go test ./... -v

gen:
	go generate ./...