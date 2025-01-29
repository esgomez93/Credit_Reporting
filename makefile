.PHONY: build run test clean docker-build docker-run

build:
	go build -o maas cmd/maas/main.go

run:
	go run cmd/maas/main.go

test:
	go test ./...

clean:
	rm -f maas

docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-down:
    docker-compose down

createdb:
    docker exec -it maas-postgres-1 psql -U $$(docker exec -it maas-postgres-1 psql -U $${POSTGRES_USER} -tAc "SELECT '\"' || current_database() || '\"'" | tr -d '"') -c "CREATE TABLE clients (client_id SERIAL PRIMARY KEY, auth_token TEXT UNIQUE NOT NULL, token_balance INTEGER DEFAULT 0)"
    docker exec -it maas-postgres-1 psql -U $$(docker exec -it maas-postgres-1 psql -U $${POSTGRES_USER} -tAc "SELECT '\"' || current_database() || '\"'" | tr -d '"') -c "CREATE TABLE api_calls (call_id SERIAL PRIMARY KEY, client_id INTEGER REFERENCES clients(client_id), timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP)"
    docker exec -it maas-postgres-1 psql -U $$(docker exec -it maas-postgres-1 psql -U $${POSTGRES_USER} -tAc "SELECT '\"' || current_database() || '\"'" | tr -d '"') -c "CREATE INDEX idx_api_calls_client_id ON api_calls (client_id)"
    docker exec -it maas-postgres-1 psql -U $$(docker exec -it maas-postgres-1 psql -U $${POSTGRES_USER} -tAc "SELECT '\"' || current_database() || '\"'" | tr -d '"') -c "INSERT INTO clients (auth_token, token_balance) VALUES ('test_token', 100)"