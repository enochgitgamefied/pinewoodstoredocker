FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .   

RUN go build -o pinewoodstore-inventory .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pinewoodstore-inventory .

# Copy static folder from builder
COPY --from=builder /app/static ./static

EXPOSE 85

CMD ["./pinewoodstore-inventory"]
