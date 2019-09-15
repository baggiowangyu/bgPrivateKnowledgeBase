# Docker-CE的安装

## CentOS 7 安装

### 移除旧版本

```
$ sudo yum remove docker \
              docker-client \
              docker-client-latest \
              docker-common \
              docker-latest \
              docker-latest-logrotate \
              docker-logrotate \
              docker-selinux \
              docker-engine-selinux \
              docker-engine
```

### 安装必要的系统工具

```
sudo yum install -y yum-utils device-mapper-persistent-data lvm2
```

### 添加软件源信息

```
sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
```

### 更新yum缓存

```
sudo yum makecache fast
```

### 安装Docker-ce

```
sudo yum -y install docker-ce
```

### 启动Docker后台服务

```
sudo systemctl start docker
```

### 测试运行hello-world

```
docker run hello-world
```
