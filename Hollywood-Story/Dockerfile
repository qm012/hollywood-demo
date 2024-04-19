FROM golang:1.21-alpine3.18

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /work
COPY . .
RUN go env && go build -ldflags="-w -s" -o server .

FROM alpine:latest

LABEL MAINTAINER="urobot@qq.com"

WORKDIR /work

COPY --from=0 /work/server ./
COPY --from=0 /work/app.yaml ./
COPY --from=0 /work/app-develop.yaml ./
COPY --from=0 /work/app-us.yaml ./
COPY --from=0 /usr/local/go/lib/time/zoneinfo.zip /opt/zoneinfo.zip

ENV ZONEINFO /opt/zoneinfo.zip

EXPOSE 8091

ENTRYPOINT ./server
