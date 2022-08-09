FROM alpine:latest
WORKDIR /app/

RUN apk add -U tzdata musl \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/shanghai" >> /etc/timezone

COPY ./wxbot /app/wxbot

ENTRYPOINT ["./wxbot"]
