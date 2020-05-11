package boot

import "github.com/gogf/gf/frame/g"

func init() {
    g.Server("Gateway").SetPort(8199)
    g.Server("Manager").SetPort(8200)
}

