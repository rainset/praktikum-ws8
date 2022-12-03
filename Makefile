migrate:
	goose -dir=./migrations sqlite3 ./nanomart.db up

run:
	go run cmd/nanomart/main.go