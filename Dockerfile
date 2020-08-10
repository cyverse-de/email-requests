FROM golang:1.14-alpine

RUN apk add --no-cache make
RUN apk add --no-cache git
RUN go get -u github.com/jstemmer/go-junit-report

ENV CGO_ENABLED=0

WORKDIR /src/email-requests
COPY . .
RUN go test ./... && \
    go build .

FROM scratch

WORKDIR /app

COPY --from=0 /src/email-requests/email-requests /bin/email-requests

ENTRYPOINT ["email-requests"]
