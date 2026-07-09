FROM golang:1.24-alpine AS builder
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码
COPY . .

# 编译
RUN RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o memo .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/memo .
COPY --from=builder /app/config ./config

# 暴露8080端口
EXPOSE 8080 
 
CMD ["./memo"]


