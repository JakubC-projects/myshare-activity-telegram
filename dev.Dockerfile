FROM golang:1.20-bullseye

RUN go install github.com/cosmtrek/air@26c752a7b07b0cdc4740a3be1e0ad84167b1a8fb

WORKDIR /app

COPY .air.toml go.mod go.sum ./
RUN go mod download

ENTRYPOINT ["air"]
