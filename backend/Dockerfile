FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend .

FROM alpine:3.24

COPY --from=builder /app/backend /backend

EXPOSE 8085 2115
ENTRYPOINT ["/backend"]
