# bgNetTunnelB 设计文档

## 模块说明

本模块作为隧道服务端，用于接受隧道客户端建立的网络连接，以此连接作为跨网络的安全通信隧道。

## 模块特性

### 隧道通信协议

隧道支持的数据承载方式：

- TCP
- UDP
- LOCAL(FTP文件摆渡)
- FTP
- SIP

两网之间通过以上方式可以进行双工通信(特别是LOCAL模式)，即可使用此隧道服务

### 隧道通信安全

- 隧道数据通信采用加密算法加密
- 通信密钥在每次构建会话时随机生成

### 隧道复用协议

隧道通信设计为单路隧道承载多路数据分发

协议结构：

| 消息头 | 协议版本号 | 协议主类型 | 协议子类型 | 映射ID | 原始客户端Address | 目标服务器Address | 网络类型 | 实际数据 |

- 消息头：固定为“GxxGm”
- 协议版本号：使用日期作为版本号，格式为“YYYYMMDD”
- 协议主类型：
  - 异常数据：“Exception”
  - 业务数据：“Business”
- 协议子类型：(暂不启用)
- 映射ID：映射端Port
- 原始客户端Address
- 目标服务器Address
- 网络类型( TCP | UDP )
- 实际数据：二进制数据流

### 隧道数据序列化方案

使用Google Protobuf v3进行数据定义与序列化

## B端服务实现思路

### B端服务的定位

- 拆解隧道数据
- 根据拆解的数据，创建客户端连接到目标服务器，并发送数据
- 客户端接收服务器返回信息，并组装数据通过隧道发回给A端
- 客户端与目标服务器连接状态发生异常，组装异常数据通过隧道发回给A端，同时销毁当前客户端

### B端服务创建客户端的思路

- 拆解数据包，拿到映射端ID(实际就是映射端Port)、原始客户端Address、目标服务器Address
- 创建一个虚拟客户端
  - 虚拟客户端启动一个协程，用于接收目标服务器返回的数据，每次目标服务器返回数据后，重新序列化数据通过隧道传给A端
  - 虚拟客户端提供发送数据接口，向目标服务器发送数据



