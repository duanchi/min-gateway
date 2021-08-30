FROM golang:alpine as go-builder
WORKDIR /go/src/github.com/duanchi/min-gateway/
COPY . /go/src/github.com/duanchi/min-gateway
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add build-base
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o heron-gateway

#FROM node:alpine as node-builder
#WORKDIR /gateway/console
#COPY ./console /gateway/console
#RUN npm config set registry "https://npm.intellidev.cn/" && \
#    npm install && npm run build

FROM alpine:latest
WORKDIR /heron
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add gettext redis nginx && \
    mkdir -p /run/nginx && \
    rm -rf /etc/nginx/conf.d/*
COPY config/application.yaml /heron/config/application.yaml
COPY config/redis.conf /etc/redis.conf
COPY console/.nginx /etc/nginx/conf.d/template/nginx.conf
COPY --from=go-builder /go/src/github.com/duanchi/min-gateway/heron-gateway /heron/gateway
#COPY --from=node-builder /gateway/console/dist /heron/console
COPY console/dist /heron/console
COPY ./docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENV CONSOLE_PORT 9800
ENV SERVER_PORT 9080
ENV AUTHORIZATION_SIGNATURE_KEY ""
ENV AUTHORIZATION_TTL 7200
ENV GATEWAY_DATA_PATH "/heron/data"
ENV GATEWAY_NATIVE_API_PREFIX "/_api"
ENV GATEWAY_NATIVE_API_ACCESS_TOKEN ""
ENV AUTHORIZATION_DSN ""
ENV ENV "development"
EXPOSE 9080
EXPOSE 9800
VOLUME ["/heron/data"]
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/heron/gateway"]