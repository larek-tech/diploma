FROM golang:1.24 AS builder

RUN apt-get update -qq 

RUN apt-get install -y -qq \
    libtesseract-dev libleptonica-dev \
    tesseract-ocr-eng tesseract-ocr-rus \
    mupdf mupdf-tools pkg-config \
    libmupdf-dev libfreetype6-dev \
    libmujs-dev libgumbo-dev libopenjp2-7-dev \
    libjbig2dec0-dev libjpeg-dev \
    libharfbuzz-dev zlib1g-dev \
    && sed -i '/^Requires: / s/$/ harfbuzz/' /usr/lib/aarch64-linux-gnu/pkgconfig/mupdf.pc \
    && echo "Updated mupdf.pc content:" \
    && cat /usr/lib/aarch64-linux-gnu/pkgconfig/mupdf.pc

ARG MODULE_NAME=github.com/larek-tech/diploma/data

WORKDIR /home/${MODULE_NAME}

COPY go.mod go.sum ./
RUN go get -t github.com/otiai10/gosseract/v2
RUN go get -t github.com/gen2brain/go-fitz
RUN go mod download

COPY . .

ENV CGO_ENABLED=1


RUN  export CGO_LDFLAGS="-lmupdf -lm -lmupdf-third -lfreetype -ljbig2dec -lharfbuzz -ljpeg -lopenjp2 -lz" \
    && go build  -o ./bin/main ./cmd/parser/main.go

FROM debian:stable-slim AS runner
ARG MODULE_NAME=github.com/larek-tech/diploma/data

WORKDIR /root/
RUN apt-get update -qq 

RUN apt-get install -y -qq \
    libtesseract-dev libleptonica-dev \
    libmupdf-dev mupdf mupdf-tools \
    libfreetype6-dev \
    libharfbuzz-dev

ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/

COPY --from=builder /home/${MODULE_NAME}/bin/main .

RUN chown root:root main

ENTRYPOINT ["/root/main"]