FROM golang:1.22.3-alpine as build-stage

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY . .

RUN go build -a -o /excange-rate-notifier main.go

FROM scratch

COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /excange-rate-notifier /excange-rate-notifier

EXPOSE 8080

ENTRYPOINT ["/excange-rate-notifier"]
