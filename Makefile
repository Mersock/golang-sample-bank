dbup:
	docker compose -f docker-compose-local.yml up -d

dbdown:
	docker compose -f docker-compose-local.yml down

createdb:
	docker exec -it golang-simple-bank-db createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it golang-simple-bank-db dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:P@ssword@localhost:5433/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:P@ssword@localhost:5433/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:P@ssword@localhost:5433/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:P@ssword@localhost:5433/simple_bank?sslmode=disable" -verbose down 1

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

kubedbup:
	kubectl apply -f k8s/db

kubedbdown:	
	kubectl delete -f k8s/db

kubeapiup:
	kubectl apply -f k8s/api  	 

kubeapidown:
	kubectl delete -f k8s/api

minikubeup:
	minikube tunnel

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: dbup createdb dropdb migrateup migrateup1 migratedown migratedown1 migratedrop sqlc test server mock kubedbup kubedbdown kubeapiup minikubeup proto

