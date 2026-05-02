FROM golang:1.25

WORKDIR /app
COPY . .

RUN go mod tidy

CMD ["go","run","main.go"]