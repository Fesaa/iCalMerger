FROM golang:1.21 as go-stage

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o /ical-merger

FROM frolvlad/alpine-glibc:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates
COPY --from=go-stage /ical-merger /app/ical-merger

EXPOSE 8080

CMD ["/app/ical-merger"]
