package device

import (
	"github.com/gogf/gf/container/gmap"
)

type DevState struct {
	Battery string
	ChargeState int
	Recording int
	Cpu string
	Storage string
	Version string
	NetDelay int
	Operator string
	Signal string
	NetType string
	WriteTime string
}

type DevGps struct {
	DeviceId string
	Latitude float32
	Longitude float32
	Height float32
	Direction float32
	Speed float32
	Radius float32
	Satellites int
	GpsAvaliable int
	GpsTime string
}

type DevInfo struct {
	Id string
	Name string
	Gps DevGps
	State DevState
	LastUpdate int64
}

type DevMgr struct {
	Devs *gmap.StrAnyMap
}

var Device_mgr DevMgr

func init() {
	Device_mgr.Devs = gmap.NewStrAnyMap(true)
}