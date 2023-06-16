default: run

run:
	go run main.go

test:
	go test ./tests -v

gen:
	go generate ./...