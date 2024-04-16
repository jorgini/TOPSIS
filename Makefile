build:
	docker compose build

run:
	docker compose up

test:
	go test -v ./app

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5432/postgres?sslmode=disable' up