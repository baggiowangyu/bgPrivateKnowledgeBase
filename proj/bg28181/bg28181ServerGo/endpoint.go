package main

import (
	"github.com/StefanKopieczek/gossip/transaction"
)

type rxInfo struct {

}

type dialog struct {

}

type endpoint struct {
	// SIP参数
	displayName	string
	username	string
	host  		string

	// 通信参数
	port		uint16
	trasport	string

	// 内部结构
	tm 			*transaction.Manager
	dialog		dialog
}
