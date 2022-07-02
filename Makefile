dbup:
	docker compose -f docker-compose-local.yml up -d

createdb:
	docker exec -it golang-simple-bank-db createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it golang-simple-bank-db dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:P@ssword@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:P@ssword@localhost:5433/simple_bank?sslmode=disable" -verbose down

migratedrop:
	migrate -path db/migration -database "postgresql://root:P@ssword@localhost:5433/simple_bank?sslmode=disable" -verbose drop

sqlc:
	sqlc generate	

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/Mersock/golang-sample-bank/db/sqlc Store

.PHONY: dbup createdb dropdb migrateup migratedown migratedrop sqlc test server mockteststore

