FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/app ./iam/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/app /app/server

COPY iam/migrations ./migrations

EXPOSE 50053

CMD ["./server"]