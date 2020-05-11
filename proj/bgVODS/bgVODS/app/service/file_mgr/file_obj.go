package file_mgr

import "os"

type VODSFileObject struct {
	Full_path	string		// 文件全路径
	File_object	*os.File	// 文件对象
}

/*
初始化文件对象
*/
func (v *VODSFileObject) Initialize(path string) error {
	var err error
	v.Full_path = path
	v.File_object, err = os.Open(v.Full_path)
	return err
}
