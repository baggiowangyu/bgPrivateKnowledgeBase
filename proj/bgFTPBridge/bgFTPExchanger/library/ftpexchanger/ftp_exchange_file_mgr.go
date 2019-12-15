package ftpexchanger

import (
	"errors"
	"strconv"
)

// 待处理的文件信息
type FileInfo struct {
	Filename	string
	Modify		int
	Size		int64
	Filetype	string
	IsMoving	bool
}

type FileMgr struct {
	// 一个map用于缓存
	Files map[string]FileInfo

}

// 查找文件记录
func (f *FileMgr) FindFileInfo(name string, modify string, size string, file_type string) (FileInfo, error) {
	key := name + modify + size + file_type
	file_info, is_exist := f.Files[key]
	if is_exist {
		return file_info, nil
	} else {
		return file_info, errors.New("Not found")
	}
}

// 增加文件记录
func (f *FileMgr) AddFile(name string, modify string, size string, file_type string) error {
	key := name + modify + size + file_type

	_, err := f.FindFileInfo(name, modify, size, file_type)
	if err != nil {
		// 说明应该是没找到对应的文件，这里构建对象添加进去
		var file_info = FileInfo{}
		file_info.Filename = name
		file_info.Filetype = file_type
		file_info.Size, err = strconv.ParseInt(size, 10, 64)
		file_info.Modify, err = strconv.Atoi(modify)

		f.Files[key] = file_info
		return nil
	} else {
		return errors.New("Exist")
	}
}

// 查找一个未摆渡的文件，放入chan中

