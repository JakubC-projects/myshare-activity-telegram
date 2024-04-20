## Build
FROM golang:1.22-bullseye as build
WORKDIR /app
COPY go.mod go.sum ./
COPY src ./src
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./src/main.go

## Deploy
FROM alpine:3.15
WORKDIR /
COPY --from=build /app/main /usr/bin/
ENTRYPOINT ["main"]
