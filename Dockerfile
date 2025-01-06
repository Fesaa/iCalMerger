# syntax=docker/dockerfile:1.3

FROM --platform=$BUILDPLATFORM golang:1.23 as go-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /ical-merger .

FROM alpine:3.21

RUN apk update && apk add --no-cache ca-certificates curl && update-ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=go-stage /ical-merger /app/ical-merger

USER 10001:10001

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "curl", "-f", "http://localhost:8080/health" ] || exit 1

LABEL maintainer="https://github.com/Fesaa/iCalMerger"
LABEL version="1.0.0"
LABEL description="iCal Merger"

CMD ["/app/ical-merger"]
