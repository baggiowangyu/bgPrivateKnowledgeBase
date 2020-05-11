package sip

/*
本文件记录了SIP协议栈的全局变量参数
*/

// 是否开启鉴权
var Sip_enable_Authenticate bool = false

// 心跳超时时间，单位：秒
var Heartbeat_check_timeout int64 = 30