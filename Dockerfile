# -----------------------------
# Stage 1: Builder
# -----------------------------
FROM golang:1.25-alpine AS builder

# Install required tools
RUN apk add --no-cache git protoc protobuf-dev

WORKDIR /app

# Install Go protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Add Go bin to PATH
ENV PATH="/root/go/bin:${PATH}"

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy project
COPY . .

# Re-generate proto files
RUN protoc \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/user.proto

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/app main.go


# -----------------------------
# Stage 2: Runtime
# -----------------------------
FROM alpine:3.19

RUN apk add --no-cache ca-certificates && update-ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /out/app .

COPY config.yaml /app/config.yaml

EXPOSE 50051

USER appuser

ENTRYPOINT ["./app"]
CMD ["serve"]
