# FIXME: use this docker image as base https://github.com/otiai10/gosseract/blob/main/Dockerfile
FROM golang:1.24 AS builder

RUN apt-get update && \
    apt-get install -y make git curl gcc g++ libtesseract-dev libleptonica-dev tesseract-ocr pkg-config \
    tesseract-ocr-eng tesseract-ocr-deu tesseract-ocr-rus

ARG MODULE_NAME=github.com/larek-tech/diploma/data

WORKDIR /home/${MODULE_NAME}

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/

RUN go build -o ./bin/main ./cmd/parser/main.go

FROM debian:stable-slim AS runner
ARG MODULE_NAME=github.com/larek-tech/diploma/data

WORKDIR /root/
RUN apt-get update && \
    apt-get install -y tesseract-ocr libtesseract-dev libleptonica-dev \
    tesseract-ocr-eng tesseract-ocr-deu tesseract-ocr-jpn && \
    rm -rf /var/lib/apt/lists/*

ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/

COPY --from=builder /home/${MODULE_NAME}/bin/main .

RUN chown root:root main

ENTRYPOINT ["/root/main"]