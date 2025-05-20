FROM scratch
WORKDIR /app
COPY exchange-log-viewer .
COPY ./public public/
ENTRYPOINT ["/app/exchange-log-viewer"]
