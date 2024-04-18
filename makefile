db-name := postgres
host := localhost
password := postgres
username := postgres
port := 5432

migrate-up:
	migrate -path migrations -database postgres://${username}:${password}@${host}:${port}/${db-name} up
migrate-down:
	migrate -path migrations -database postgres://${username}:${password}@${host}:${port}/${db-name} down
run:
	go mod download
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main/main.go
	./main

