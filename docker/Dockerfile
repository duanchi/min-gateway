FROM golang:alpine as go-builder
WORKDIR /go/src/github.com/duanchi/min-gateway/
COPY . /go/src/github.com/duanchi/min-gateway
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add build-base
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    CGO_ENABLED=0 GOOS=linux go build -a -o min-gateway

FROM node:alpine as node-builder
WORKDIR /gateway/console
COPY ./console /gateway/console
RUN npm config set registry "https://npm.intellidev.cn/" && \
    npm install && npm run build

FROM alpine:latest
WORKDIR /min
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add gettext redis nginx && \
    mkdir -p /run/nginx && \
    rm -rf /etc/nginx/http.d/*
COPY config/application.yaml /min/config/application.yaml
COPY config/redis.conf /etc/redis.conf
#COPY ui/.nginx /etc/nginx/http.d/template/nginx.conf
#COPY gateway /min/gateway
COPY --from=go-builder /go/src/github.com/duanchi/min-gateway/min-gateway /min/gateway
COPY --from=node-builder /gateway/console/dist /min/console
COPY ui/dist /min/console
COPY ./docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENV CONSOLE_PORT 9800
ENV SERVER_PORT 9080
ENV AUTHORIZATION_SIGNATURE_KEY ""
ENV AUTHORIZATION_TTL 7200
ENV GATEWAY_DATA_PATH "/min/data"
ENV GATEWAY_CONSOLE_API_PREFIX "/_api"
ENV GATEWAY_CONSOLE_API_ACCESS_TOKEN ""
ENV GATEWAY_NATIVE_API_PREFIX "/native"
ENV GATEWAY_CONFIG_REFRESH_INTERVAL 30
ENV AUTHORIZATION_DEFAULT_SIGNLETON "false"
ENV AUTHORIZATION_DSN ""
ENV ENV "production"
EXPOSE 9080
EXPOSE 9800
VOLUME ["/min/data"]
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/min/gateway"]