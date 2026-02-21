FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY app-bin /usr/local/bin/

EXPOSE 5002

CMD ["app-bin"]
