build:
	go build -o /bin/server ./cmd/tg_bot/main.go

run: build
	@./bin/server

test:
	go test -v ./...

image:
	docker build -t myserver .

compose-up: image
	docker-compose -f ./docker-compose.yml up -d

compose-down:
	docker-compose -f docker-compose.yml down -v