FROM golang:1.21 as builder

WORKDIR /app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

# cache
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -ldflags '-w -s' -o e5bot .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && apk add --no-cache ca-certificates

RUN mkdir build && cp e5bot build && mv config.yml.example build/config.yml

FROM alpine:latest

RUN apk add tzdata
COPY --from=builder /app/build /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/e5bot", "-c", "config.yaml"]
