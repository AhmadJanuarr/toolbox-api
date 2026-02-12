# Menggunakan image dasar Golang versi terbaru (Alpine linux biar ringan)
FROM golang:1.24-alpine

# Set working directory di dalam container
WORKDIR /app

# Install dependency sistem yang dibutuhkan untuk library gambar (CGO)
# gcc, musl-dev, dll sering dibutuhkan untuk kompilasi library gambar yang pakai C
RUN apk add --no-cache gcc musl-dev

# Copy file go.mod dan go.sum terlebih dahulu (biar cache layer optimal)
COPY go.mod go.sum ./

# Download semua dependency Go
RUN go mod download

# Copy semua source code project ke dalam container
COPY . .

# Build aplikasi Go menjadi binary bernama 'main'
RUN go build -o main .

# Expose port yang digunakan aplikasi (misal 8080)
EXPOSE 8080

# Perintah yang dijalankan saat container start
CMD ["./main"]