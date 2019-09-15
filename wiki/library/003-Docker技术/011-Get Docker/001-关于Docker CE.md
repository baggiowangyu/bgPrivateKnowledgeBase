# 关于Docker CE

原文地址：https://docs.docker.com/install/

Docker Community Edition (CE)是个体开发人员和小型团队开始使用Docker和尝试基于容器的应用程序较为理想的一种解决方案。

Docker CE有三种更新通道：

- Stable(稳定版)：给我们提供最新的可用版本；
- Test(测试版)：提供稳定版之前的一个预览版；
- Nightly(每日版/开发版)：提供为下一个发布版本为目标的正在开发的版本；

## Releases

对于Docker CE引擎

## Support(支持)

Docker CE 在第一次year-month分支发布后的每7个月进行一次版本发布。
Docker EE 在第一次year-month分支发布后的每24个月进行一次版本发布。

这意味着BUG报告以及分支移植在全生命周期都需要被评估。

在year-month分支到达了生命周期重点，这个分支才允许从代码仓库中移除。

### Reporting security issues(报告安全问题)

Docker的维护人员都在认真的对待安全问题。如果你发现了一个安全问题，请马上让维护人员注意！

请勿公开此类问题；请将相关问题发邮件给security@docker.com。

非常感谢您的安全报告，Docker将为此公开感谢您。Docker也喜欢送礼物——如果你喜欢swag，一定要让我们知道。Docker目前没有提供付费的安全奖励计划，但不排除未来的可能性。

### Supported platforms

Docker CE可用于多个平台。使用以下表格为您选择最佳安装路径。

#### DESKTOP(桌面端)

| Platform(平台) | x86_64 |
|:-:|:-:|
| Docker for Mac(macOS) | √ |
| Docker for Windows(Microsoft Windows 10) | √ |

#### SERVER(服务器端)

| Platform(平台) | x86_64 / amd64 | ARM | ARM64 / AARCH64 | IBM Power(ppc64le) | IBM Z(s390x) |
|:-:|:-:|:-:|:-:|:-:|:-:|
| CentOS | √ | - | √ | - | - |
| Debian | √ | √ | √ | - | - |
| Fedora | √ | - | √ | - | - |
| Ubuntu | √ | √ | √ | √ | √ |

### Backporting(移植)

Docker公司将对Docker产品的Backports进行优先级排序。Docker员工或代码仓库维护者将努力确保合理的错误修复使之成为活动的版本。

如果有一些重要的修复应该被考虑，以使backport成为活动的发布分支，那么一定要在PR描述中强调这一点，或者在PR中添加注释。


### Upgrade path

补丁版本总是向后兼容其年月版本。

## Not covered

一般来说，本文档中未提及的任何内容都可能在任何版本中更改。

## Exceptions(异常)

异常的产生是由于安全补丁引起的。如果需要中断发布过程或产品功能，它将被清楚地传达，解决方案将被考虑在总体影响之下。
