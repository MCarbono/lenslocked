default: run

run:
	go run main.go

test:
	go test ./controllers -v