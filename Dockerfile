# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o sommelier .

# Run stage
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/sommelier .
COPY --from=builder /app/frontend/dist ./frontend/dist

# /data will be the mounted volume for file-based DB (create with fly volumes create data --size 1)
ENV PORT=8080

EXPOSE 8080

CMD ["./sommelier"]
