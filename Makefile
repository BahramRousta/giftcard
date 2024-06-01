build:
	go build -o giftCrad

migrate:
	go run main.go migrate

run:
	go run main.go

postgres:
	docker run -d \
		--name postgres \
		-e POSTGRES_USER=root \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=gift_card_db \
		-e POSTGRES_SSL_MODE=disable \
		-e TZ=Asia/Tehran \
		-p 5433:5432 \
		postgres:16.3-alpine3.20

redis:
	docker run -d --name redis -p 6379:6379 redis:latest

.PHONY: run migrate build