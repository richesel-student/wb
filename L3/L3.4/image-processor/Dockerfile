FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-app ./cmd/api
RUN go build -o worker-app ./cmd/worker


FROM alpine:3.18
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/api-app .
COPY --from=builder /app/worker-app .


COPY --from=builder /app/internal/handler/static ./static

CMD ["./api-app"]