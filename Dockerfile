FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ls -la /app/docs

RUN go build -o /app/main ./cmd/main.go

RUN chmod +x /app/main

CMD ["/app/main"]
