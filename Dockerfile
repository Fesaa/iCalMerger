FROM golang:1.21 as go-stage

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o /ical-merger

FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates
COPY --from=go-stage /ical-merger /app/ical-merger

EXPOSE 8080

CMD ["/app/ical-merger"]
