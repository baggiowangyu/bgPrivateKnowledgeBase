package data

type GxxGmLocationInfo struct {
	Usercode 	string		// 用户编码
	Timevalue	string		// 时间（UTC时间？）
	Longitude	string		// 经度，正值为东经，负值为西经
	Latitude	string		// 纬度，正值为北纬，负值为南纬
	Speed		string		// GPS速度信息
	Direction 	string		// 方向（一般需要依陀螺仪）
	Altitude	string		// 海拔高度
	Radius		string		// 定位半径，精度
}

type GxxGmExceptionInfo struct {
	Usercode 	string		// 用户编码
	Storage		string		// 存储模块异常
	Battery		string		// 电池模块异常
	Ccd			string		// 拍摄模块异常
	Mic			string		// 麦克风异常
	Position	string		// 设备静止异常
}

type GxxGmDeviceBaseInfo struct {
	Usercode 	string		// 用户编码
	Carrieroperator	string	// 无线网络运行商
	Nettype		string		// 无线网络类型
	Signal		string		// 无线信号强度
	Battery		string		// 剩余电池电量
	Storage		string		// 剩余存储空间
	Cpu			string		// CPU占用率
	Version		string		// 前端设备版本号
	LocalRecord	string		// 前端设备本地录像状态
	ChargeState	string		// 前端设备充电状态
}
