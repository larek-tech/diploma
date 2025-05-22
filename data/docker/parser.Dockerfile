FROM golang:1.24 AS builder

WORKDIR /app

# Install dependencies in a single layer to reduce image size
# Explicitly install libjpeg62-turbo-dev to ensure version 62
RUN apt-get update -qq && \
    apt-get install -y --no-install-recommends \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr-eng \
    tesseract-ocr-rus \
    libjpeg62-turbo-dev \
    mupdf-tools \
    libmupdf-dev \
    libfreetype6-dev \
    libopenjp2-7-dev \
    libjbig2dec0-dev \
    libharfbuzz-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set Tessdata path
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/

# Copy and build the Go application
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/main ./cmd/parser/main.go

# Runtime stage
FROM debian:stable-slim

WORKDIR /root/

# Install ONLY the necessary runtime libraries
RUN apt-get update -qq && \
    # First install essential packages
    apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    libjpeg62-turbo && \
    # Check what mupdf packages are available
    apt-cache search mupdf | grep -E '^lib' || true && \
    # Now install other dependencies
    apt-get install -y --no-install-recommends \
    libtesseract5 \
    liblept5 \
    tesseract-ocr-eng \
    tesseract-ocr-rus && \
    # Try to install libmupdf-dev, but continue if it fails
    apt-get install -y --no-install-recommends libmupdf-dev || \
    # Install additional dependencies if libmupdf-dev isn't available
    apt-get install -y --no-install-recommends \
    libfreetype6 \
    libopenjp2-7 \
    libjbig2dec0 \
    libharfbuzz0b && \
    # Create a symlink to ensure libjpeg.so points to version 62
    if [ -f /usr/lib/x86_64-linux-gnu/libjpeg.so.62 ]; then \
    rm -f /usr/lib/x86_64-linux-gnu/libjpeg.so; \
    ln -sf /usr/lib/x86_64-linux-gnu/libjpeg.so.62 /usr/lib/x86_64-linux-gnu/libjpeg.so; \
    fi && \
    # Download and add server certificate for s3.diploma.larek.tech
    echo "Downloading certificate for s3.diploma.larek.tech..." && \
    mkdir -p /usr/local/share/ca-certificates && \
    update-ca-certificates && \
    # Verify libjpeg configuration
    ls -la /usr/lib/x86_64-linux-gnu/libjpeg* && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set environment variables
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/
# Force the dynamic loader to use the version 62 library
ENV LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjpeg.so.62
# Trust self-signed certificates if needed for development
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

# Copy the compiled binary from builder
COPY --from=builder /app/bin/main .
RUN chmod +x main

ENTRYPOINT ["/root/main"]