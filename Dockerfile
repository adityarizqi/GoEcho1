FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-krs-app main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /go-krs-app .

COPY templates ./templates
COPY .env .

EXPOSE 8080

CMD ["./go-krs-app"]