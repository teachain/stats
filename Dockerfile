FROM golang:1.23.2-alpine AS builder

RUN  set -eux \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache gcc musl-dev linux-headers git  \
    && rm -vrf /var/cache/apk/


ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN  go version && go env &&  \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64  go build \
    -ldflags "-linkmode 'external' -extldflags '-static' -X 'main.gitCommit=$(git rev-parse HEAD)' -X 'main.buildDate=$(date +%Y-%m-%d)' -X 'main.gitDate=$(git show -s --format=%cd $(git rev-parse HEAD) --date=short)'" \
    -v -o stats ./main.go


FROM alpine:latest
#  设置上海时区环境变量
ENV TZ=Asia/Shanghai
# # 更换apk源,安装时区以及curl依赖，设置系统时区为上海时区，如需安装依赖只需在curl后添加依赖名称即可
RUN  set -eux \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache  ttf-dejavu fontconfig tzdata  curl busybox-extras iputils ca-certificates \
    && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && rm -vrf /var/cache/apk/

WORKDIR /app/logs

WORKDIR /app

###将编译产物拷贝至镜像
COPY --from=builder  /app/stats  /app/stats

