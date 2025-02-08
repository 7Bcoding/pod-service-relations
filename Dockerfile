FROM golang:latest
LABEL maintainer="cenquanyu@xunlei.com"

# RUN apt-get update && apk add --no-cache tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# set working directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

ADD ./output /usr/src/app/
ENTRYPOINT  ["/usr/src/app/bin/pod-service-relations"]

EXPOSE 8080