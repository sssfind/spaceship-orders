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
COPY order/ ./order/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/app ./order/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/app /app/server

COPY --from=builder /app/order/migrations ./migrations

EXPOSE 8080

CMD ["./server"]