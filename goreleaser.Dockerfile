FROM scratch
WORKDIR /app
COPY sync-groups-from-authentik-to-outline .
ENTRYPOINT ["/app/sync-groups-from-authentik-to-outline"]
