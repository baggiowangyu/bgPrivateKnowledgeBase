package service

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gmutex"
)

type FaceInfo struct {
	DeviceId string
	Device_latitude float32
	Device_longtitude float32
	Img_path string
}

type CarNumberInfo struct {
	DeviceId string
	Device_latitude float32
	Device_longtitude float32
	Car_number string
	Img_path string
}

type DataSyncSerice struct {
	face_infos []FaceInfo
	face_lock *gmutex.Mutex

	car_number_infos []CarNumberInfo
	car_number_lock *gmutex.Mutex
}

var DataSync DataSyncSerice

func (d *DataSyncSerice) Initialize() error {
	// 定义一个定时任务
	d.face_lock = gmutex.New()
	_, err := gcron.Add("30 * * * * *", d.SyncFaceInfo)
	if err != nil {
		glog.Error(err)
		return err
	}

	d.car_number_lock = gmutex.New()
	_, err = gcron.Add("30 * * * * *", d.SyncCarNumberInfo)
	if err != nil {
		glog.Error(err)
		return err
	}

	gcron.Start("")
	return err
}

func (d *DataSyncSerice) AddFaceInfo(info FaceInfo)  {
	//
	d.face_lock.Lock()
	d.face_infos = append(d.face_infos, info)
	d.face_lock.Unlock()
}

func (d *DataSyncSerice) AddCarNumberInfo(info CarNumberInfo) {
	//
	d.car_number_lock.Lock()
	d.car_number_infos = append(d.car_number_infos, info)
	d.car_number_lock.Unlock()
}

func (d *DataSyncSerice) SyncFaceInfo()  {
	var err error

	d.face_lock.Lock()
	defer d.face_lock.Unlock()

	tx, err := g.DB().Begin()
	if err != nil {
		glog.Error(err)
		return
	}

	for _, info := range d.face_infos {
		_, err := tx.Insert("face_record", gdb.Map{
			"device_id"         : info.DeviceId,
			"device_latitude"   : info.Device_latitude,
			"device_longtitude" : info.Device_longtitude,
			"img_path"          : info.Img_path,
		})

		if err != nil {
			glog.Error(err)
			//回滚
			_ = tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		glog.Error(err)
		_ = tx.Rollback()
		return
	}

	// 清空数组
	for index := 0; index < len(d.face_infos); index++ {
		d.face_infos = append(d.face_infos[:index], d.face_infos[index+1:]...)
	}
}

func (d *DataSyncSerice) SyncCarNumberInfo()  {
	var err error

	d.car_number_lock.Lock()
	defer d.car_number_lock.Unlock()

	tx, err := g.DB().Begin()
	if err != nil {
		glog.Error(err)
		return
	}

	for _, info := range d.car_number_infos {
		_, err := tx.Insert("car_record", gdb.Map{
			"device_id"         : info.DeviceId,
			"device_latitude"   : info.Device_latitude,
			"device_longtitude" : info.Device_longtitude,
			"car_number"        : info.Car_number,
			"img_path"          : info.Img_path,
		})

		if err != nil {
			glog.Error(err)
			//回滚
			_ = tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		glog.Error(err)
		_ = tx.Rollback()
		return
	}

	// 清空数组
	for index := 0; index < len(d.car_number_infos); index++ {
		d.car_number_infos = append(d.car_number_infos[:index], d.car_number_infos[index+1:]...)
	}
}