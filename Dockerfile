FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY web/ web/
COPY db/migrations/ db/migrations/
EXPOSE 8088
CMD ["./server"]
