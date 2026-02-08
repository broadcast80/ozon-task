FROM golang:1.24.2-alpine AS build-stage

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CG0_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]