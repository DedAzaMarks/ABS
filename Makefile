build:
	go build -o /bin/server ./cmd/tg_bot/main.go

run: build
	@/bin/server

test:
	go test -v ./...

#image:
#	docker build -t server .

compose-build:
	docker-compose -f ./docker-compose.yml build

up: compose-build
	docker-compose -f ./docker-compose.yml up

down:
	docker-compose -f docker-compose.yml down -v