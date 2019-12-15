package bridgemgr

import (
	"github.com/gogf/gf/net/ghttp"
)

func init() {

}

type BridgeMgrController struct {

}

// 添加一个桥接数据
func (b *BridgeMgrController) AddBridge(r *ghttp.Request) {
	// 这里主要是向数据库添加一条记录
}

// 移除一个桥接数据
func (b *BridgeMgrController) RemoveBridge(r *ghttp.Request) {

}

// 查询所有桥接数据
func (b *BridgeMgrController) QueryBridge(r *ghttp.Request) {

}
