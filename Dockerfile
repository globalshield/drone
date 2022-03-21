FROM golang:1.17-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/github.com/globalshield/drone

COPY . .

RUN go mod download && go mod verify

RUN go build -v -o /go/bin/scanner .

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/scanner /app/scanner

ENTRYPOINT ["/app/scanner"]
