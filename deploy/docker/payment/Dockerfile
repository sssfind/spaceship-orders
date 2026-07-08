FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY go.mod go.sum ./
COPY shared/go.mod shared/go.sum* ./shared/
COPY inventory/go.mod inventory/go.sum* ./inventory/
COPY order/go.mod order/go.sum* ./order/
COPY payment/go.mod payment/go.sum* ./payment/

RUN go mod download

COPY shared/ ./shared/
COPY payment/ ./payment/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/app ./payment/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/app /app/server

EXPOSE 50052

CMD ["./server"]