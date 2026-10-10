FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o gateway .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && \
    adduser -S app -G app

COPY --from=builder /app/gateway .

EXPOSE 8080

USER app

CMD ["./gateway"]
