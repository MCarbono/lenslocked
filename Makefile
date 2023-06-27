default: run

run:
	go run main.go

test:
	go test ./tests/./... -v

unit_test:
	go test ./tests/unit -v

integration_test:
	go test ./tests/integration -v
	
api_test:
	go test ./tests/api -v

gen:
	go generate ./...