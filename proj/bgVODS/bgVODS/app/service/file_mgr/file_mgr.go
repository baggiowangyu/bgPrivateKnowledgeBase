/*
文件管理模块：

此模块需要任务如下：
1、管理服务所用根目录

*/

package file_mgr

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"time"
)

type VODSFileMgr struct {
	root_path string
}

var VodsFileMgr VODSFileMgr

/*
初始化文件管理器
*/
func (v *VODSFileMgr) Initialize(root string) error {
	v.root_path = root
	return nil
}

func (v *VODSFileMgr) IfFileExists(relative_path string) bool {
	real_path := v.root_path + relative_path

	_, err := os.Stat(real_path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func (v *VODSFileMgr) GetFileSize(relative_path string) (int64, error) {
	real_path := v.root_path + relative_path

	fileInfo, err := os.Stat(real_path)
	if os.IsNotExist(err) {
		return 0, err
	}

	// 获取文件大小
	filesize := fileInfo.Size()
	return filesize, err
}

func (v *VODSFileMgr) GetFileModifyTimeString(relative_path string) (string, error) {
	real_path := v.root_path + relative_path

	fileInfo, err := os.Stat(real_path)
	if os.IsNotExist(err) {
		return "", err
	}

	// 获取文件最后修改时间
	filemodifytime := fileInfo.ModTime().Format(time.RFC1123)
	return filemodifytime, err
}

func (v *VODSFileMgr) GetFileETag(relative_path string) (string, error) {
	real_path := v.root_path + relative_path

	real_file_object, err := os.Open(real_path)
	if err != nil {
		return "", err
	}

	// 获取文件的ETag
	md5_obj := md5.New()
	_, _ = io.Copy(md5_obj, real_file_object)
	etag := hex.EncodeToString(md5_obj.Sum(nil))
	_ = real_file_object.Close()

	return etag, err
}

func (v *VODSFileMgr) OpenFile(relative_path string) (*VODSFileObject, error) {
	real_path := v.root_path + relative_path
	file_object := new(VODSFileObject)
	err := file_object.Initialize(real_path)
	if err != nil {
		return nil, err
	} else {
		return file_object, err
	}
}
