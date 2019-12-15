package router

import (
	"bgFTPBridgeServer/app/api/bridgemgr"
	"github.com/gogf/gf/frame/g"
)

func init() {
	bridge_mgr_controller := new(bridgemgr.BridgeMgrController)
	g.Server().BindObject("POST:/BridgeMgr", bridge_mgr_controller, "AddBridge")
	g.Server().BindObject("POST:/BridgeMgr", bridge_mgr_controller, "RemoveBridge")
	g.Server().BindObject("GET:/BridgeMgr", bridge_mgr_controller, "QueryBridge")
	
}