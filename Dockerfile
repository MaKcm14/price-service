FROM golang:1.23-alpine3.21

WORKDIR /price-service

COPY . .

VOLUME /price-service/logs

RUN go mod tidy

WORKDIR /price-service/scripts/linux

# DEBUG port: for checking its work
EXPOSE 8080

ENTRYPOINT [ "go", "run", "../../cmd/app/main.go" ]
