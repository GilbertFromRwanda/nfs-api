FROM golang:1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app-bin .


FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app-bin /usr/local/bin/

EXPOSE 5002

CMD ["app-bin"]
