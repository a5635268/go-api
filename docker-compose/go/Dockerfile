FROM golang:alpine as builder
# 创建镜像1，用于编译go
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src
RUN apk add wget \
    && wget https://github.com/a5635268/go-api/archive/master.zip \
    && unzip master.zip \
    && cd /go/src/go-api-master \
    && go env \
    && go list \
    && go build -v -a -o go-api .

# 二段编译： 创建mini镜像2，用于运行go
# docker build . -t  go-api
# docker run -p 9001:8000 go-api-master
FROM scratch
COPY --from=builder /go/src/go-api-master/config /config
COPY --from=builder /go/src/go-api-master/go-api /
#CMD ["/bin/bash"]
ENTRYPOINT ["/go-api","server","-c=config/settings.dev.yml"]
