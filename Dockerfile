FROM golang:1.21 as go-stage

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /ical-merger

FROM frolvlad/alpine-glibc:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=go-stage /ical-merger /app/ical-merger

CMD ["/app/ical-merger"]
