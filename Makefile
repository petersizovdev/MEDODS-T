build:
	@go build -o bin/build cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/build

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run migration/migration.go up

migrate-down:
	@go run migration/migration.go down