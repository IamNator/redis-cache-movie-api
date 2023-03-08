# Step 1: Modules caching
FROM golang:latest


WORKDIR /app

COPY . .

RUN go mod tidy

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

CMD ["go", "test", "-v", "./integration-test/..."]