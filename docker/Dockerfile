ARG BASE=registry.cn-shanghai.aliyuncs.com/winc-edgex/edgex:go1.23.1-alpine3.16
FROM ${BASE} AS builder


# set the working directory
WORKDIR /simple-driver

COPY . .

RUN #GOCGO=CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o  ./cmd/build/driver-linux-amd64  ./cmd


FROM registry.cn-shanghai.aliyuncs.com/winc-edgex/alpine:3.13

COPY --from=builder /simple-driver/cmd/build/driver-linux-amd64 /app/

RUN apk add --no-cache tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone \
&& apk del tzdata

EXPOSE 1883

RUN mkdir /logs


# Start
CMD [ "/bin/sh", "-c", "/app/driver-linux-amd64" ]


