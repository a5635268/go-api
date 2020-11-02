FROM golang:alpine as builder
# 创建镜像1，用于编译go
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/go-api
COPY . .
RUN go env && go list && go build -v -a -o go-api .

# 二段编译： 创建mini镜像2，用于运行go
# docker build . -t  go-api
# docker run -p 9001:8000 go-api
FROM scratch
COPY --from=builder /go/src/go-api/config /config
COPY --from=builder /go/src/go-api/go-api /
#CMD ["/bin/bash"]
ENTRYPOINT ["/go-api","server","-c=config/settings.dev.yml"]
