.PHONY: all postgres createdb dropdb migrateup migratedown sqlc test format peepdb migration_test server proto
DB_NAME=
DB_URI=postgresql://root:password@localhost:5432/qwetu_petro_db?sslmode=disable
all: test

postgres:
	docker run --name datapoint_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine

createdb:
	docker exec -it datapoint_db createdb --username=root --owner=root $(DB_NAME)

dropdb:
	docker exec -it datapoint_db dropdb $(DB_NAME)

migratedown:
	migrate -path db/migrations -database ${DB_URI} -verbose down 1

prodmigrateup:
	migrate -path db/migrations -database ${PROD_DB_URI} -verbose up

sqlc:
	sqlc generate

test:
	@echo "Running tests..." && \
	go test -v ./... 2>&1 | tee test_output.log ; \
	test_status=$$? ; \
	echo "Tests complete"; \
	exit $$test_status

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/qwetu_petro/backend/db/sqlc Store

format:
	go fmt ./... -w

peepdb:
	docker exec -it datapoint_db psql --username=root --dbname=$(DB_NAME)

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

migration_test: dropdb createdb

server:
	go run main.go

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger  --openapiv2_opt=allow_merge=true,merge_file_name=qwetu_backend_api  \
	proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl


redis:
	docker run --name redis -p 6379:6379 -d redis:7.2-rc2