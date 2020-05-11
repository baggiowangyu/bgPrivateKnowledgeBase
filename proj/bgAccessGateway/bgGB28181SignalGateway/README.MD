# bgGB28181SignalGateway

本项目遵循《GB/T 28181-2016 公共安全视频监控联网系统信息传输、交换、控制技术要求》内容进行开发。

## 项目组件

- 程序框架：GoFrame 1.11.4
- SIP协议栈：github.com/stefankopieczek/gossip v0.0.0-20161014194501-bf923f8025bc

## 程序说明

### 程序结构定义

- 在boot.go中引入SIP协议服务```import _ "bgGB28181SignalGateway/app/service/sip"```
- 项目路径```bgGB28181SignalGateway/app/service/sip```中实现了本程序主体：SIP信令服务
- 项目路径```bgGB28181SignalGateway/app/service/data```中定义了本程序中的数据结构

### 程序流程

SIP服务器的入口在```bgGB28181SignalGateway/app/service/sip/sip_server.go```中的**```init()```**函数内。

#### sip_server.go

文件内定义了一个结构体：

```cassandraql
/*
SIP服务端对象
持有一个sip客户端清单
持有一个配置参数：是否开启鉴权
持有一个SIP传输管理器
持有一个SIP事务处理器
*/
type SipServer struct {
	// SIP客户端清单
	Sip_clients map[string]*SipClient
	Clients_map_lock *gmlock.Locker

	// SIP传输管理器
	Sip_transport_manager transport.Manager
	// SIP事务处理器
	Sip_transaction_manager *transaction.Manager

	// 鉴权开关
	Sip_enable_Authenticate bool
}
```

结构体内有5个成员：

- Sip_clients：已接入接入客户端清单；
- Clients_map_lock：客户端map读写锁；
- Sip_transport_manager：SIP协议栈传输管理器；
- Sip_transaction_manager：SIP协议栈事务处理器；
- Sip_enable_Authenticate：SIP鉴权开关；

##### init()

1. 创建一个SipServer结构对象，调用启动函数StartUp()

##### SipServer.StartUp()

1. 初始化客户端map清单；
2. 初始化客户端map读写锁；
3. 初始化SIP协议栈传输管理器，目前写死了是基于**```UDP```**协议；
4. 初始化SIP协议栈事务管理器，初始化监听地址写死为**```0.0.0.0:5060```**；
5. 配置服务鉴权开关，默认为**```关（False）```**；
6. 启动一个协程，开始接收SIP消息；

##### SipServer.SipRequestRecv()

1. 
