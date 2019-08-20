FROM golang:1.12.9-alpine3.10 as golang

WORKDIR /go/src/app
COPY . .

RUN apk add --no-cache --update git &&\
    go get -d -v ./... &&\
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . &&\
    rm -rf /var/cache/apk/*

FROM alpine:3.10

RUN mkdir /opt/speedlog
WORKDIR /opt/speedlog
COPY --from=0 /go/src/app/main .

EXPOSE 8012
