FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser -D -u 10001 appuser
WORKDIR /app

COPY --from=builder /app/main ./main

EXPOSE 8080

USER appuser

CMD ["./main"]
