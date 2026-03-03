FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
COPY static ./static
RUN go build -o registry-ui .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/registry-ui .
EXPOSE 8080
CMD ["./registry-ui"]
