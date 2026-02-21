ARG GO_VERSION=1

FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify

RUN go build -v -o /app-bin .


FROM debian:bookworm

COPY --from=builder /app-bin /usr/local/bin/

EXPOSE 8080

CMD ["app-bin"]