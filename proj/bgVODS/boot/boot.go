package boot

import "github.com/gogf/gf/frame/g"

func init() {
    g.Server().SetPort(g.Config().GetInt("app.http_port"))
}

