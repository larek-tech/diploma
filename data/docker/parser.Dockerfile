# Builder
FROM golang:1.24-alpine AS builder
RUN apk add --update make git curl

ARG MODULE_NAME=github.com/larek-tech/diploma/data

RUN apk upgrade --update-cache --available && \
    apk add openssl && \
    rm -rf /var/cache/apk/*
COPY Makefile /home/${MODULE_NAME}/Makefile
COPY go.mod /home/${MODULE_NAME}/go.mod
COPY go.sum /home/${MODULE_NAME}/go.sum

WORKDIR /home/${MODULE_NAME}

# Only download dependencies if go.mod or go.sum changed
RUN go mod download

# Now copy the rest of the source code
COPY . /home/${MODULE_NAME}

RUN CGO_ENABLED=0 go build -o ./bin/main ./cmd/parser/main.go

# Service
FROM alpine:latest AS runner
ARG MODULE_NAME=github.com/larek-tech/diploma/data

WORKDIR /root/

COPY --from=builder /home/${MODULE_NAME}/bin/main .

RUN chown root:root main

ENTRYPOINT ["/root/main"]