FROM golang:latest

RUN mkdir "/app"
WORKDIR /app
COPY . .

RUN go mod tidy

CMD ["go", "run", "cmd/server/main.go"]