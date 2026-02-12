# Stage 1: Build Stage
FROM golang:1.24-alpine AS builder

# Install C Compiler (gcc, musl-dev) yang WAJIB untuk library webp
RUN apk add --no-cache build-base

WORKDIR /app

# Copy dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi
# CGO_ENABLED=1 wajib untuk library chai2010/webp
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 2: Runtime Stage (Image Akhir)
FROM alpine:latest

WORKDIR /root/

# Copy binary dari builder
COPY --from=builder /app/main .

# Buat folder temp yang dibutuhkan aplikasi
RUN mkdir -p temp/uploads temp/processed temp/compressed temp/resized

# Expose port
EXPOSE 8080

# Jalankan
CMD ["./main"]