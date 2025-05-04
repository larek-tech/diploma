# Builder
FROM golang:1.24-alpine AS builder

RUN apk update && \
    apk upgrade --update-cache --available && \
    apk add --no-cache make git curl gcc musl-dev openssl librdkafka-dev

ARG MODULE_NAME=github.com/larek-tech/diploma/data

WORKDIR /home/${MODULE_NAME}

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1

RUN go build -tags musl -o ./bin/main ./cmd/crawler/main.go

# Runner
FROM alpine:latest AS runner
ARG MODULE_NAME=github.com/larek-tech/diploma/data

WORKDIR /root/

COPY --from=builder /home/${MODULE_NAME}/bin/main .

RUN chown root:root main

ENTRYPOINT ["/root/main"]