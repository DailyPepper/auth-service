FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY migrations/ ./migrations/

COPY . .
RUN ls -la migrations/

RUN go build -o auth-service cmd/auth/main.go

EXPOSE 50051

CMD ["./auth-service"]