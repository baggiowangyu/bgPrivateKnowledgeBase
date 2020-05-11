package dao

type WorkStationInfo struct {

}

type WorkStationHeartBeat struct {
	Ws_gbcode 				string
	Ws_ip 					string
	Ws_online 				int
	Ws_enable 				int
	Ws_total_storage_size 	int64
	Ws_used_storage_size 	int64
	Ws_cpu_usage 			int
	Ws_ram_usage 			int
	Ws_raid_state 			int
	Ws_alarm_type 			int
	Ws_version 				string
}