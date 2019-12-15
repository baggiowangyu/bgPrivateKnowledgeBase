package bridgemgr

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/gfile"
)

// 声明一个全局变量
var Bridge_Object = new(BridgeObject)


// 桥接处理函数
func bridge_handler(conn *gtcp.Conn)  {
	defer conn.Close()
	// 首先检查该连接是否已经缓存下来
	client_socket_name := conn.Conn.RemoteAddr().String()
	is_conn_exist := Bridge_Object.IfConnExists(client_socket_name)

	// 连接不存在，则创建连接
	if !is_conn_exist {
		Bridge_Object.SaveConnRecord(client_socket_name, conn)
	}

	for {
		// 接收数据，接收多少就传多少，摆渡文件命名：[sock_name]_[current_time].exchange
		data, err := conn.Recv(-1)
		if len(data) > 0 {
			// 接收数据，写入发送端
			file_path := g.Config()
			gfile.Create()
		}

		if err != nil {
			break
		}
	}
}

// 定义一个桥接对象
type BridgeObject struct {
	Client_Conns map[string]interface{}
}

// 启动桥接对象
func (bo *BridgeObject) StartUp() error {
	g.TCPServer().SetHandler(bridge_handler)
	err := g.TCPServer().Run()
	return err
}

// 停止桥接对象
func (bo *BridgeObject) Stop() error {
	err := g.TCPServer().Close()
	return err
}

// 查看连接对象是否存在
func (bo *BridgeObject) IfConnExists(sock_name string) bool {
	for socket_name, _ := range bo.Client_Conns {
		if sock_name == socket_name {
			return true
		}
	}

	return false
}

// 保存连接对象
func (bo *BridgeObject) SaveConnRecord(sock_name string, conn_object *gtcp.Conn) {
	bo.Client_Conns[sock_name] = conn_object
}

// 移除连接对象
func (bo *BridgeObject) RemoveConnRecord(sock_name string)  {
	delete(bo.Client_Conns, sock_name)
}

