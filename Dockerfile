FROM golang:1.21 as go-stage

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o /ical-merger

FROM debian:stable-slim

WORKDIR /app

COPY --from=go-stage /ical-merger /app/ical-merger

EXPOSE 8080

CMD ["/app/ical-merger"]
