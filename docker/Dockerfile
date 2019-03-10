# Build binary
FROM golang:1.11-alpine as builder
ENV GODIR=/go/src/golang-proxy-server
ENV GO111MODULE=on
COPY . $GODIR
# touch a file for build code
WORKDIR $GODIR
RUN apk --update add --no-cache git

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o golang-proxy-server ./main.go

FROM alpine
ENV USER=golang-proxy-server TERM=xterm
ENV HOME=/home/$USER
WORKDIR /

# Setting golang-proxy-server 
RUN mkdir /etc/golang-proxy-server /var/log/golang-proxy-server
COPY --from=builder /go/src/golang-proxy-server/golang-proxy-server /usr/bin/golang-proxy-server
RUN apk add --no-cache tzdata ca-certificates \
    && rm -rf /var/cache/apk/* 

# Change TimeZone
RUN apk add --update tzdata
ENV TZ=Asia/Taipei
# Clean APK cache
RUN rm -rf /var/cache/apk/*

# add config 
ADD ./config /etc/golang-proxy-server

ADD ./docker/docker-entrypoint.sh /entrypoint.sh

ENTRYPOINT ["sh", "/entrypoint.sh"]