migrate:
	goose -dir=./migrations sqlite3 ./nanomart.db up

run:
	go run cmd/nanomart/main.go

test:
	go test -count=1 ./...

test-integration:
	go test -tags integration -count=1 ./...