package subscribe

import "github.com/gogf/gf/container/gmap"

type SubscribeInfo struct {
	Id string
	Method int
	Url string
}

type SubscribeMgr struct {
	SubscribeInfos gmap.StrAnyMap
}

func (s *SubscribeMgr) Add(info *SubscribeInfo) error {
	internal_info, is_find := s.SubscribeInfos.Search(info.Id)
	if !is_find {
		// 增加新的
		s.SubscribeInfos.Set(info.Id, info)
	} else {
		real_info := internal_info.(*SubscribeInfo)
		real_info.Method = info.Method
		real_info.Url = info.Url
		//s.SubscribeInfos.
	}
}