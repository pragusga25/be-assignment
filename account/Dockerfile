# Gunakan image golang official sebagai base image
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod dan sum files
COPY go.mod go.sum ./

# Download semua dependencies
RUN go mod download

# Copy source code dari current directory ke working directory di dalam container
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Mulai stage baru untuk mendapatkan image yang lebih kecil
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy executable dari builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080
EXPOSE 8080

# Command untuk menjalankan executable
CMD ["./main"]
