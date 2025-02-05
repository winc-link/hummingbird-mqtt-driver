ARG BASE=registry.cn-shanghai.aliyuncs.com/winc-edgex/edgex:1.18.8-alpine3.16
FROM ${BASE} AS builder


# set the working directory
WORKDIR /simple-driver

COPY . .

RUN go build -ldflags "-s -w" -o  ./cmd/build/driver-linux-amd64  ./cmd

FROM registry.cn-shanghai.aliyuncs.com/winc-edgex/alpine:3.12


RUN apk add tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" >  /etc/timezone \
&& apk del tzdata \


COPY --from=builder /simple-driver/cmd/build/driver-linux-amd64 /app/

EXPOSE 1883

RUN mkdir /logs


# Start
CMD [ "/bin/sh", "-c", "/app/driver-linux-amd64" ]


