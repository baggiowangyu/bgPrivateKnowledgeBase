package router

import (
	"../app/api/orgmgr"
	"github.com/gogf/gf/frame/g"
)

func init() {
	g.Server().BindObject("/orgmgr", new(orgmgr.Controller))
}
