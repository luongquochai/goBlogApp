DB_URL=postgresql://root:secret@localhost:5432/myblogdb?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

start_docker:
	docker start postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root myblogdb

dropdb:
	docker exec -it postgres dropdb myblogdb

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/main.go

.PHONY: postgres start_docker createdb dropdb migrateup migratedown sqlc test server