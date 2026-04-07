FROM golang:1.25 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o app ./cmd

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/web ./web  

CMD ["./app"]