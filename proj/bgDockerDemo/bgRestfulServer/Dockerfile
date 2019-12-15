# 基础镜像：CentOS 7.0
FROM centos:7

# 指定镜像维护者
MAINTAINER baggiowangyu@gmail.com

# 尝试将程序拷贝到指定程序目录
# 路径一定要写对，不然会报拷贝错误
COPY ./bgRestfulService/ ./baggio_app/

# 创建一个数据挂载点
# 从本地主机或者其他容器挂载数据卷，一般用来存放数据库和需要保存的数据等
#VOLUME ["/data"]

# 如果使用了dubbo-go，这里需要先指定环境变量
#ENV CONF_PROVIDER_FILE_PATH XXXXX
#ENV APP_LOG_CONF_FILE XXXXX

# 启动服务
#RUN ["/bin/bash", "-c", "./app/bgRestfulServer/bgRestfulServer"]

### 文件编辑结束后，在终端执行docker build命令来创建镜像。
# docker build -t baggiowangyu/bg-restful-service .
### 通过镜像启动容器
# docker run -p 9009:9009 -t -i baggiowangyu/bg-restful-service:latest ./baggio_app/bgRestfulServer