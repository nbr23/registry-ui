FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETOS TARGETARCH
WORKDIR /app
COPY go.mod .
COPY main.go .
COPY static ./static
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o registry-ui .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/registry-ui .
ENV PORT=8080
EXPOSE ${PORT}
CMD ["./registry-ui"]
