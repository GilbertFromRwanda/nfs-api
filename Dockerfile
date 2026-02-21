FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY app-bin /usr/local/bin/

EXPOSE 5002

CMD ["app-bin"]
