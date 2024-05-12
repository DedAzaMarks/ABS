FROM golang:latest AS build-dependencies
WORKDIR /go/src/app
COPY ["go.mod", "go.sum", "/go/src/app/"]
RUN go mod download -x && go mod vendor

FROM golang:latest AS build
WORKDIR /go/src/app
COPY --from=build-dependencies ["/go/src/app/vendor", "/go/src/app/"]
COPY ["go.mod", "go.sum", "Makefile", "/go/src/app/"]
COPY ["cmd", "/go/src/app/cmd"]
COPY ["internal", "/go/src/app/internal"]
ENV CGO_ENABLED=0
RUN ["make", "build"]

FROM ubuntu AS certificates
RUN apt-get update && apt-get install -y ca-certificates

FROM scratch as server
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build ["/bin/server", "/bin/"]
COPY [".env", "/"]
CMD ["/bin/server"]