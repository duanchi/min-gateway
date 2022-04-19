FROM docker.in-mes.io/library/golang as builder
WORKDIR /go/src/github.com/duanchi/min-gateway/
COPY . /go/src/github.com/duanchi/min-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /min-gateway


FROM alpine:latest
WORKDIR /min
COPY --from=builder /min-gateway /min/gateway
#COPY ./bin/min-gateway /min/gateway
COPY config/application.yaml /min/config/application.yaml
COPY ./docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh && chmod +x /min/gateway
ENV SERVER_PORT 9080
ENV AUTHORIZATION_SIGNATURE_KEY ""
ENV AUTHORIZATION_TTL 7200
ENV GATEWAY_DATA_PATH "/min/data"
ENV GATEWAY_CONSOLE_API_PREFIX "/_api"
ENV GATEWAY_CONSOLE_API_ACCESS_TOKEN ""
ENV GATEWAY_NATIVE_API_PREFIX "/native"
ENV CONFIG_REFRESH_INTERVAL 30
ENV AUTHORIZATION_DEFAULT_SIGNLETON "false"
ENV ENABLE_GRAY_INSTANCE "true"
ENV AUTHORIZATION_DSN ""
ENV DB_DSN ""
ENV CACHE_DSN ""
ENV ENV "production"
EXPOSE 9080
VOLUME ["/min/data"]
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/min/gateway"]