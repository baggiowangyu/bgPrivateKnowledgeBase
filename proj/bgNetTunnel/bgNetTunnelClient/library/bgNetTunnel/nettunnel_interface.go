package bgNetTunnel

// 定义一个回调函数
type NetTunnelServerRecvCallback func(date []byte) error

// 这里定义一个抽象类，定义网络通道接口
// 这里的抽象规则其实比较简单，这里不关心我们的通道协议层面的内容
// 这里只对两端网络进行抽象，例如：将FTP抽象为接口，将TCP抽象为接口，将UDP抽象为接口等
type NetTunnelServerInterface interface {

	// 发送数据
	SendData(data []byte) error

	// 接收数据，这里传递一个参数到
	RecvData(callback NetTunnelServerRecvCallback)
}
