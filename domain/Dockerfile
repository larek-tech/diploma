FROM golang:1.24.2-alpine AS builder
WORKDIR /builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/server/main.go

FROM scratch AS prod
WORKDIR /prod

COPY config/config.yaml ./config/
COPY --from=builder /builder/main .

CMD [ "./main" ]