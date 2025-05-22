FROM golang:1.24 AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update -qq && \
    apt-get install -y --no-install-recommends \
    libtesseract-dev \
    libleptonica-dev \
    tesseract-ocr-eng \
    tesseract-ocr-rus \
    libjpeg62-turbo-dev \
    poppler-utils && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Match builder env with runtime env
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/

# Build application
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/main ./cmd/parser/main.go

# Runtime stage
FROM debian:stable-slim

WORKDIR /root/

# Install runtime dependencies (single RUN to reduce layers)
RUN apt-get update -qq && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    libjpeg62-turbo \
    libtesseract5 \
    liblept5 \
    tesseract-ocr-eng \
    tesseract-ocr-rus \
    poppler-utils \
    libfreetype6 \
    libopenjp2-7 \
    libjbig2dec0 \
    libharfbuzz0b && \
    if [ -f /usr/lib/x86_64-linux-gnu/libjpeg.so.62 ]; then \
    rm -f /usr/lib/x86_64-linux-gnu/libjpeg.so; \
    ln -sf /usr/lib/x86_64-linux-gnu/libjpeg.so.62 /usr/lib/x86_64-linux-gnu/libjpeg.so; \
    fi && \
    mkdir -p /usr/local/share/ca-certificates && \
    update-ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Environment variables
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/
ENV LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjpeg.so.62
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

# Copy compiled binary from builder
COPY --from=builder /app/bin/main .
RUN chmod +x main

# Verify tesseract data files exist
RUN ls -la $TESSDATA_PREFIX

ENTRYPOINT ["/root/main"]