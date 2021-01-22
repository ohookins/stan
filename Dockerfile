# Build stage
FROM golang:alpine AS builder
WORKDIR /build
COPY *.go /build/
RUN go build -o stan

# Run stage
FROM golang:alpine
WORKDIR /app
COPY --from=builder /build/stan .
EXPOSE 8080
ENTRYPOINT ["/app/stan"]