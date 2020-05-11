package device

import (
	"GxxGmGoDeviceInfoPusher/app/service/bgBase"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
)

var (
	// HTTP推送配置
	http_enable = g.Config().GetInt("http_push.enable")
	http_push_url = g.Config().GetString("http_push.url")
	http_method = g.Config().GetString("http_push.method")

	// TCP推送配置
	tcp_enable = g.Config().GetInt("tcp_push.enable")
	tcp_push_url = g.Config().GetString("tcp_push.host")
	tcp_method = g.Config().GetInt("tcp_push.port")

	// Kafka推送配置
	kafka_enable = g.Config().GetInt("kafka_push.enable")
	kafka_push_url = g.Config().GetString("kafka_push.host")
	kafka_method = g.Config().GetInt("kafka_push.port")
	kafka_topic = g.Config().GetString("kafka_push.topic")
)

func init() {
	// 启动一个定时任务，每秒钟将Redis数据同步到本地缓存
	gcron.Add("@every 1s", SyncDeviceInfoTask)

	// 启动第二个定时任务，根据配置的推送协议与推送地址，向指定目标推送信息
	gcron.Add("@every 5s", PushDeviceInfoTask)

	// 启动第三个定时任务，根据配置检查所有设备信息列表，当每次检查时间与上次检查时间间隔超过阈值，则将设备信息清理掉，并发出设备下线通知
	gcron.Add("@every 10s", OnlineDeviceCheckTask)
}

func SyncDeviceInfoTask() {
	// 同步设备信息任务
	glog.Debug("Run SyncDeviceInfoTask()")

	// 首先，尝试拿到所有在线设备ID
	res, err := g.Redis().Do("KEYS", "online.status.*")
	if err != nil {
		glog.Error(err)
		return
	}

	for _, element := range res.([]interface{}) {
		res_string := bgBase.Uint8Array_2_String(element.([]uint8))
		device_id := res_string[14:]
		glog.Debugf("SyncDeviceInfoTask() >>> 拿到设备[%s]的在线状态。", device_id)

		// 拿Status信息
		var status_json *gjson.Json
		param := fmt.Sprintf("status.%s", device_id)
		res, err = g.Redis().Do("GET", param)
		if err != nil {
			glog.Error(err)
		} else {
			if res != nil {
				status_string := bgBase.Uint8Array_2_String(res.([]uint8))
				status_json = gjson.New(status_string)
			}
		}



		// 拿GPS信息
		var gps_json *gjson.Json
		param = fmt.Sprintf("pos_%s", device_id)
		res, err = g.Redis().Do("GET", param)
		if err != nil {
			glog.Error(err)
		} else {
			if res != nil {

				gps_string := bgBase.Uint8Array_2_String(res.([]uint8))
				gps_json = gjson.New(gps_string)
			}
		}




		// 当前时间
		current_timestamp := gtime.Now().Unix()

		// 首先尝试寻找缓冲中是否存在此设备
		device_info_interface, is_find := Device_mgr.Devs.Search(device_id)
		if !is_find {
			// 没找到，则创建新的
			device_info_interface := &DevInfo{}
			device_info_interface.LastUpdate = current_timestamp
			device_info_interface.Id = device_id

			if status_json != nil {
				device_info_interface.State.Battery = status_json.GetString("battery")
				device_info_interface.State.Storage = status_json.GetString("storage")
				device_info_interface.State.ChargeState = status_json.GetInt("chargeState")
				device_info_interface.State.Recording = status_json.GetInt("recording")
				device_info_interface.State.Cpu = status_json.GetString("cpu")
				device_info_interface.State.Version = status_json.GetString("version")
				device_info_interface.State.NetDelay = status_json.GetInt("netDelay")
				device_info_interface.State.Operator = status_json.GetString("operator")
				device_info_interface.State.Signal = status_json.GetString("signal")
				device_info_interface.State.NetType = status_json.GetString("nettype")
				device_info_interface.State.WriteTime = status_json.GetString("writeTime")
			}

			if gps_json != nil {
				device_info_interface.Gps.GpsTime = gps_json.GetString("gpstime")
				device_info_interface.Gps.DeviceId = gps_json.GetString("deviceId")
				device_info_interface.Gps.Latitude = gps_json.GetFloat32("latitude")
				device_info_interface.Gps.Longitude = gps_json.GetFloat32("longitude")
				device_info_interface.Gps.Height = gps_json.GetFloat32("height")
				device_info_interface.Gps.Direction = gps_json.GetFloat32("direction")
				device_info_interface.Gps.Speed = gps_json.GetFloat32("speed")
				device_info_interface.Gps.Radius = gps_json.GetFloat32("radius")
				device_info_interface.Gps.Satellites = gps_json.GetInt("satellites")
				device_info_interface.Gps.GpsAvaliable = gps_json.GetInt("gpsavailable")
			}

			Device_mgr.Devs.Set(device_info_interface.Id, device_info_interface)

		} else {
			// 找到了，更新字段，首先更新时间字段
			device_info := device_info_interface.(*DevInfo)
			device_info.LastUpdate = current_timestamp
			//device_info.Id string
			//device_info.Name string

			// State，首先检查Redis内缓存的与内存缓存的写入时间是否相同，若相同，则不更新，否则更新
			if status_json != nil {
				if device_info.State.WriteTime != status_json.GetString("writeTime") {
					device_info.State.Battery = status_json.GetString("battery")
					device_info.State.Storage = status_json.GetString("storage")
					device_info.State.ChargeState = status_json.GetInt("chargeState")
					device_info.State.Recording = status_json.GetInt("recording")
					device_info.State.Cpu = status_json.GetString("cpu")
					device_info.State.Version = status_json.GetString("version")
					device_info.State.NetDelay = status_json.GetInt("netDelay")
					device_info.State.Operator = status_json.GetString("operator")
					device_info.State.Signal = status_json.GetString("signal")
					device_info.State.NetType = status_json.GetString("nettype")
					device_info.State.WriteTime = status_json.GetString("writeTime")
				}
			}


			// Gps，如果卫星时间发生变化，则更新
			//device_info.Gps
			if gps_json != nil {
				if device_info.Gps.GpsTime != gps_json.GetString("gpstime") {
					device_info.Gps.GpsTime = gps_json.GetString("gpstime")
					device_info.Gps.DeviceId = gps_json.GetString("deviceId")
					device_info.Gps.Latitude = gps_json.GetFloat32("latitude")
					device_info.Gps.Longitude = gps_json.GetFloat32("longitude")
					device_info.Gps.Height = gps_json.GetFloat32("height")
					device_info.Gps.Direction = gps_json.GetFloat32("direction")
					device_info.Gps.Speed = gps_json.GetFloat32("speed")
					device_info.Gps.Radius = gps_json.GetFloat32("radius")
					device_info.Gps.Satellites = gps_json.GetInt("satellites")
					device_info.Gps.GpsAvaliable = gps_json.GetInt("gpsavailable")
				}
			}
		}
	}
}

func PushDeviceInfoTask() {
	// 推送设备信息任务
	glog.Debug("Run PushDeviceInfoTask()")

	// 遍历所有设备
	Gpss := make([]DevGps, 0)

	Device_mgr.Devs.Iterator(func(k string, v interface{}) bool {
		devinfo := v.(*DevInfo)

		// 只有GPS数据可用时加入
		if devinfo.Gps.GpsAvaliable == 1{
			Gpss = append(Gpss, devinfo.Gps)
		}
		return true
	})

	GpsData := make(map[string]interface{}, 0)
	GpsData["total"] = len(Gpss)
	GpsData["data"] = Gpss

	GpsDataJson := gjson.New(GpsData)
	GpsDataJsonString := GpsDataJson.Export()

	glog.Debugf("待推送的定位信息：%s", GpsDataJsonString)

	if http_enable == 1 {
		// 向指定的http地址推送定位信息
		if http_method == "post" {
			response, err := ghttp.Post(http_push_url, GpsDataJsonString)
			if err != nil {
				glog.Error(err)
			}

			glog.Debugf("收到推送结果：%d", response.StatusCode)
		}

	}
}

func OnlineDeviceCheckTask() {
	// 设备在线状态检查
	glog.Debug("Run OnlineDeviceCheckTask()")

	Device_mgr.Devs.Iterator(func(k string, v interface{}) bool {
		devinfo := v.(*DevInfo)

		// 比较状态信息最后更新时间与当前时间差，若超出15秒，则踢下线
		current_time := gtime.Now().Unix()
		if current_time - devinfo.LastUpdate > 15000 {
			Device_mgr.Devs.Remove(k)

			// 通知下线
		}
		return true
	})
}
