# Get Docker CE for CentOS

原文地址：https://docs.docker.com/install/linux/docker-ce/centos

要在CentOS上开始使用Docker CE，请确保满足先决条件，然后安装Docker。

## Prerequisites(先决条件)

### Docker EE customers(Docker EE客户)

要安装Docker Enterprise Edition (Docker EE)，请转到[CentOS的Docker EE](https://docs.docker.com/install/linux/docker-ee/centos/, "CentOS的Docker EE")，而不是本主题。

有关Docker EE的更多信息，请参见[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/, "Docker Enterprise Edition")。

### OS requirements(操作系统要求)

要安装Docker CE，您需要一个CentOS 7的维护版本。不支持或测试存档版本。

必须启用 **```centos-extras```** 仓库。默认情况下，这个存储库是启用的，但是如果您禁用了它，则需要重新启用它。

建议使用 **```overlay2```** 存储驱动程序。

### Uninstall old versions(卸载旧版本)

Docker的旧版本称为 **```Docker```** 或 **```Docker-engine```** 。如果安装了这些组件，请卸载它们以及相关的依赖项。

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

如果 **```yum```** 报告这些包都没有安装，则没有问题。

保存 **```/var/lib/docker/```** 的内容，包括图像、容器、卷和网络。Docker CE包现在称为 **```Docker-CE```** 。

## Install Docker CE(安装Docker CE)

根据需要，可以以不同的方式安装Docker CE：

- 大多数用户设置Docker的代码仓库并从中进行安装，以简化安装和升级任务。这是推荐的方法。
- 有些用户下载RPM包并手动安装，并完全手动管理升级。这在一些情况下非常有用，比如在没有互联网接入的物理隔离系统上安装Docker。
- 在测试和开发环境中，一些用户选择使用自动化的方便脚本来安装Docker。

### Install using the repository(使用代码仓库安装)

在新主机上首次安装Docker CE之前，需要设置Docker代码仓库。之后，您可以从代码仓库中安装和更新Docker。

#### SET UP THE REPOSITORY(设置代码仓库)

1. 安装所需要的包。**```yum-utils```** 提供了 **```yum-config-manager```** 功能，并且 **```devicemapper```** 存储驱动修安排 **```device-mapper-persistent-data```** 和 **```lvm2```**。

  ```
      $ sudo yum install -y yum-utils \
             device-mapper-persistent-data \
             lvm2
  ```

2. 使用以下命令来设置稳定版代码仓库。即使我们要安装edeg或者test代码仓库，稳定版的代码仓库也是一定需要安装的。

  ```
      $ sudo yum-config-manager \
          --add-repo \
          https://download.docker.com/linux/centos/docker-ce.repo
  ```

3. 可选项：启动edeg和测试代码仓库。这些代码仓库包含在 **```docker.repo```** 文件中，但是默认情况下是被禁用的。我们可以在stable代码仓库中启用它。

  ```
      $ sudo yum-config-manager --enable docker-ce-edge
  ```

  ```
      $ sudo yum-config-manager --enable docker-ce-test
  ```

  我们可以运行 **```yum-config-manager```** 命令，携带参数 **```--disable```** 来关闭edge和test代码仓库。要重新开启，则携带 **```--enable```** 参数。
  下面的命令展示了如何禁用edge代码仓库：

  ```
      $ sudo yum-config-manager --disable docker-ce-edge
  ```

#### INSTALL DOCKER CE

1. 安装最新版本的Docker CE，或进入下一步安装特定版本:

  ```
      $ sudo yum install docker-ce
  ```

  如果提示接受GPG密钥，验证指纹与 **```060A 61C5 1B55 8A7F 742B 77AA C52F EB6B 621E 9F35```** 匹配，如果匹配则接受。

  ```
  有多个Docker代码仓库的情况：
  如果我们有多个代码仓库是可用的
  ```








































-
