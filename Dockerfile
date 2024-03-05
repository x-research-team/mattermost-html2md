FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o service.exe ./cmd/server/cmd/main.go

EXPOSE 8080

CMD ["./service.exe"]
