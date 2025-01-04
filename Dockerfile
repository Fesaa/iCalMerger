# syntax=docker/dockerfile:1.3

FROM --platform=$BUILDPLATFORM golang:1.23 as go-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /ical-merger .

FROM scratch

WORKDIR /app

COPY --from=go-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-stage /ical-merger /app/ical-merger

# Non root user and group
USER 10001:10001

# Healthcheck
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "curl", "-f", "http://localhost:8080/health" ] || exit 1

LABEL maintainer="https://github.cowm/Fesaa/iCalMerger"
LABEL version="1.0.0"
LABEL description="iCal Merger"

CMD ["/app/ical-merger"]
