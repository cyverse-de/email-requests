FROM golang:1.21

ENV CGO_ENABLED=0

WORKDIR /src/email-requests
COPY . .
RUN go test ./... && \
    go build .

FROM debian:stable-slim

WORKDIR /app

COPY --from=0 /src/email-requests/email-requests /bin/email-requests

ENTRYPOINT ["email-requests"]
