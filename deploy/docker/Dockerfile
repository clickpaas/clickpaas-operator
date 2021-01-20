FROM alpine  as builder
ARG VERSION="v8"
ENV GOPROXY https://goproxy.cn
RUN apk update && apk add git go && mkdir /root/go/src -pv && \
    cd /root/go/src/ && \
    https://github.com/clickpaas/clickpaas-operator.git && \
    cd clickpaas-operator && go mod vendor && go build .


FROM alpine
WORKDIR /root
COPY --from=builder /root/go/src/clickpaas-operator/clickpaas-operator .
ENTRYPOINT ["./clickpaas-operator"]