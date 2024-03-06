FROM golang:1.21-alpine

RUN apk add --no-cache git

ENV CGO_ENABLED=0

WORKDIR /src/email-requests
COPY . .
RUN go test ./... && \
    go build .

FROM scratch

WORKDIR /app

COPY --from=0 /src/email-requests/email-requests /bin/email-requests

ENTRYPOINT ["email-requests"]
