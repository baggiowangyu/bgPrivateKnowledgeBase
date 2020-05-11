package sip

import "github.com/stefankopieczek/gossip/transaction"

/*
此文件定义了端点信息
以及端点的
*/

type endpoint struct {
	// Sip Params
	displayName string
	username    string
	host        string

	// Transport Params
	port      uint16 // Listens on this port.
	transport string // Sends using this transport. ("tcp" or "udp")

	// Internal guts
	tm       *transaction.Manager
	dialog   dialog
	dialogIx int
}

type dialog struct {
	callId    string
	to_tag    string // The tag in the To header.
	from_tag  string // The tag in the From header.
	currentTx txInfo // The current transaction.
	cseq      uint32
}

type txInfo struct {
	tx     transaction.Transaction // The underlying transaction.
	branch string                  // The via branch.
}