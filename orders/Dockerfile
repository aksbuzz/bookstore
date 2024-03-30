FROM golang:1.20-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]