package boot

import (
	"bgBaseKernel/app/api/consul"
)

func init() {
	go consul.InitializeConsul()
}

