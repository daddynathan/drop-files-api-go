FROM golang:1.25.1-alpine AS builder

RUN apk update && apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o dropfiles ./cmd/server/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/dropfiles .
COPY --from=builder /app/.env .
COPY --from=builder /app/uploads ./uploads

RUN mkdir -p ./uploads && chmod -R 775 ./uploads

CMD ["./dropfiles"]