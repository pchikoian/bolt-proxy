FROM golang:1.24 as builder

WORKDIR /go/src/app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o bolt-proxy .

FROM gcr.io/distroless/base:latest
COPY --from=builder /go/src/app/bolt-proxy /usr/local/bin/bolt-proxy

ENV BOLT_PROXY_BIND=0.0.0.0:7687
EXPOSE 7687/tcp

ENTRYPOINT ["/usr/local/bin/bolt-proxy"]
