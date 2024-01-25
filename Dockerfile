# Dockerfile
# 指定基础镜像
FROM golang:latest

# 设置工作目录
WORKDIR /app

# 拷贝项目代码到容器中
COPY . .

# 编译可执行文件
RUN go build -o main .

# 设置执行命令
CMD ["./main"]