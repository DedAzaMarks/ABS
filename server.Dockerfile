FROM golang:latest AS build

WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -o /bin/server ./cmd/server/main.go

FROM scratch
COPY --from=build /bin/server /bin/server
CMD ["/bin/server"]