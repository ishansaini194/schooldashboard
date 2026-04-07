FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o school ./cmd

FROM alpine:latest
WORKDIR /app

RUN mkdir -p /app/data
COPY --from=builder /app/school .

EXPOSE 8080
CMD ["./school"]
