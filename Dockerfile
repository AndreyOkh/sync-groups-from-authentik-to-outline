FROM golang:1.24-alpine AS build
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o exchange-log-viewer ./

FROM scratch
LABEL authors="okhotnikov"
WORKDIR /app
COPY --from=build ["/build/sync-groups-from-authentik-to-outline", "."]
ENTRYPOINT ["/app/sync-groups-from-authentik-to-outline"]
