/*

客户端管理服务：

由于客户端管理的是连接到此服务器的所有客户端信息，其生命周期与此服务相同
因此，相关信息全部缓存于内存中，不做持久化处理

 */

package bgClientMgr

import (
	"errors"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/net/gudp"
)

type ClientObject struct {
	Client_addr 	string
	Mapping_addr 	string
	Tcp_conn		*gtcp.Conn
	Udp_conn		*gudp.Conn
	Conn_type		string
}

type ClientMgr struct {
	Client_list map[string]*ClientObject
}

func (c *ClientMgr) Initialize() {
	c.Client_list = make(map[string]*ClientObject, 0)
}

// 查找客户端是否存在
func (c *ClientMgr) FindClientObject(cli_addr string) (*ClientObject, error) {
	object, exist := c.Client_list[cli_addr]
	if exist {
		return object, nil
	} else {
		return object, errors.New("Client object not found.")
	}
}

// 增加一个已连接进来的客户端信息
func (c *ClientMgr) AddClientObject(cli_addr string, mapping_addr string, conn_type string, tcp_conn *gtcp.Conn, udp_conn *gudp.Conn) {

	client_obj := new(ClientObject)
	client_obj.Client_addr = cli_addr
	client_obj.Mapping_addr = mapping_addr
	client_obj.Conn_type = conn_type
	if client_obj.Conn_type == "TCP" {
		client_obj.Tcp_conn = tcp_conn
	} else {
		client_obj.Udp_conn = udp_conn
	}

	c.Client_list[cli_addr] = client_obj
}

// 移除一个已断开的连接
func (c *ClientMgr) RemoveClientObject(cli_addr string) {
	delete(c.Client_list, cli_addr)
}
