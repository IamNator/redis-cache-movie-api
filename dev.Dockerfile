FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go install github.com/githubnemo/CompileDaemon@latest

CMD CompileDaemon -build="go build -o api cmd/api/main.go" -command="./api"

