package vods

import (
	"bgVODS/boot"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func split(s rune) bool {
	if s == '-' {
		return true
	}
	return false
}

func OpenFile(path string) (*os.File, error) {
	real_file_object, err := os.Open(path)
	if err != nil {
		glog.Error(err)
	}
	return real_file_object, err
}

func PlayHandler(r *ghttp.Request) {
	// 拼装文件全路径
	real_path := boot.Global_vods_root + r.RequestURI

	real_file_object, err := OpenFile(real_path)
	if err != nil {
		r.Response.WriteStatus(500)
		return
	}

	block_size := int64(4096 * 100)
	read_buf := make([]byte, block_size)

	for {
		readed_bytes, err := real_file_object.Read(read_buf)
		if err != nil {
			//glog.Error(err)
			break
		}

		r.Response.Write(read_buf[:readed_bytes])
	}

	_ = real_file_object.Close()
}

func PlayHandlerEx(r *ghttp.Request) {
	// 拼装文件全路径
	real_path := boot.Global_vods_root + r.RequestURI

	// 检查文件是否存在
	fileInfo, err := os.Stat(real_path)
	if os.IsNotExist(err) {
		glog.Infof("Request %s failed. File not found.", r.RequestURI)
		r.Response.WriteStatus(404)
		return
	}

	// 获取文件大小
	filesize := fileInfo.Size()
	// 获取文件最后修改时间
	filemodifytime := fileInfo.ModTime()

	real_file_object, err := OpenFile(real_path)
	if err != nil {
		r.Response.WriteStatus(500)
		return
	}

	// 获取文件的ETag
	md5_obj := md5.New()
	_, _ = io.Copy(md5_obj, real_file_object)
	etag := hex.EncodeToString(md5_obj.Sum(nil))
	_ = real_file_object.Close()

	// 然后检查Header中是否包含Range
	var range_exist bool = false
	var range_string string
	var start, end, length int64
	for key, val := range r.Request.Header {
		if strings.EqualFold("Range", key) {
			// 找到了偏移信息，然后来分析具体的偏移范围
			// Range: bytes=start-end, 这里我们需要对字符串进行切割，得到start和end值
			range_string = val[0][6:]
			glog.Infof("[VODS::Controller::PlayHandler] Request header key %s, value : %s", key, range_string)

			// 假如
			results := strings.FieldsFunc(range_string, split)
			if len(results) < 2 {
				start, err = strconv.ParseInt(results[0], 10, 64)
				if err != nil {
					glog.Error(err)
					r.Response.WriteStatus(500)
					return
				}
				end = filesize
			} else {
				// if len(results) == 2
				start, err = strconv.ParseInt(results[0], 10, 64)
				if err != nil {
					glog.Error(err)
					r.Response.WriteStatus(500)
					return
				}

				end, err = strconv.ParseInt(results[1], 10, 64)
				if err != nil {
					glog.Error(err)
					r.Response.WriteStatus(500)
					return
				}
			}

			length = end - start
			range_exist = true
			break
		}
	}

	/*
		返回数据，这里比较讲究。
		我们在返回的时候需要在响应头加上“Accept-Ranges : bytes”，表示本服务器支持断点续传
		- 没有在请求头中携带“Range”的，我们直接返回所有数据，包含：
		  - 【响应头】200 OK
		  - 【响应头】Accept-Ranges : bytes
	      - 【响应头】ETag : "45b9ef0-5911062c10d4c"
		  - 【响应头】Last-Modified: Tue, 27 Aug 2019 02:54:05 GMT
		  - 【响应头】Content-Type : video/mp4
		  - 【响应头】Content-Length : [需要传输的总长度]
		  - 【响应体】文件数据
		- 请求头中携带“Range”的，说明是进行断点续传了，我们返回的数据应当有所变化，包含：
		  - 【响应头】206 Partial Content
		  - 【响应头】Accept-Ranges : bytes
		  - 【响应头】ETag : "45b9ef0-5911062c10d4c"
		  - 【响应头】Last-Modified : Tue, 27 Aug 2019 02:54:05 GMT
		  - 【响应头】Content-Range : bytes [start]-[end]/[total]
		  - 【响应头】Content-Type : video/mp4
		  - 【响应头】Content-Length : [需要传输的总长度]
		  - 【响应体】文件数据
	*/

	// 通用响应头
	response_header := r.Response.Header()
	response_header["Accept-Ranges"] = append(response_header["Accept-Ranges"], "bytes")
	response_header["ETag"] = append(response_header["ETag"], etag)
	response_header["Last-Modified"] = append(response_header["Last-Modified"], filemodifytime.Format(time.RFC1123))

	// 文件传输块
	block_size := int64(4096 * 100)
	read_buf := make([]byte, block_size)

	real_file_object, err = OpenFile(real_path)
	if err != nil {
		r.Response.WriteStatus(500)
		return
	}

	if range_exist {
		// 断点续传，准备一下响应头
		Content_Range := fmt.Sprintf("bytes %d-%d/%d", start, end, filesize)
		response_header["Content-Range"] = append(response_header["Content-Range"], Content_Range)
		r.Response.WriteHeader(206)

		// 读取数据，发送出去
		_, err := real_file_object.Seek(start, 0)
		if err != nil {
			glog.Error(err)
			r.Response.WriteStatus(500)
			_ = real_file_object.Close()
			return
		}

		// 这里有可能导致申请失败
		read_buf = make([]byte, length)
		readed_bytes, err := real_file_object.Read(read_buf)
		if err != nil {
			glog.Error(err)
			r.Response.WriteStatus(500)
			_ = real_file_object.Close()
			return
		}

		r.Response.Write(read_buf[:readed_bytes])


	} else {
		// 读取数据，发送出去
		for {
			readed_bytes, err := real_file_object.Read(read_buf)
			if err != nil {
				//glog.Error(err)
				break
			}

			r.Response.Write(read_buf[:readed_bytes])
		}
	}

	_ = real_file_object.Close()
}

func PlayHandlerSecurity(r *ghttp.Request) {
	// 拼装文件全路径
	real_path := boot.Global_vods_root + r.RequestURI

	// 检查文件是否存在
	fileInfo, err := os.Stat(real_path)
	if os.IsNotExist(err) {
		glog.Infof("Request %s failed. File not found.", r.RequestURI)
		r.Response.WriteStatus(404)
		return
	}

	// 获取文件大小
	filesize := fileInfo.Size()
	// 获取文件最后修改时间
	filemodifytime := fileInfo.ModTime()

	real_file_object, err := OpenFile(real_path)
	if err != nil {
		r.Response.WriteStatus(500)
		return
	}

	// 获取文件的ETag
	md5_obj := md5.New()
	_, _ = io.Copy(md5_obj, real_file_object)
	etag := hex.EncodeToString(md5_obj.Sum(nil))
	_ = real_file_object.Close()

	// 然后检查Header中是否包含Range
	var range_exist bool = false
	var range_string string
	var start, end, length int64
	for key, val := range r.Request.Header {
		if strings.EqualFold("Range", key) {
			// 找到了偏移信息，然后来分析具体的偏移范围
			// Range: bytes=start-end, 这里我们需要对字符串进行切割，得到start和end值
			range_string = val[0][6:]
			glog.Infof("[VODS::Controller::PlayHandler] Request header key %s, value : %s", key, range_string)

			// 假如
			results := strings.FieldsFunc(range_string, split)
			if len(results) < 2 {
				start, err = strconv.ParseInt(results[0], 10, 64)
				if err != nil {
					glog.Error(err)
					r.Response.WriteStatus(500)
					return
				}
				end = filesize
			} else {
				// if len(results) == 2
				start, err = strconv.ParseInt(results[0], 10, 64)
				if err != nil {
					glog.Error(err)
					r.Response.WriteStatus(500)
					return
				}

				end, err = strconv.ParseInt(results[1], 10, 64)
				if err != nil {
					glog.Error(err)
					r.Response.WriteStatus(500)
					return
				}
			}

			length = end - start
			range_exist = true
			break
		}
	}

	/*
		返回数据，这里比较讲究。
		我们在返回的时候需要在响应头加上“Accept-Ranges : bytes”，表示本服务器支持断点续传
		- 没有在请求头中携带“Range”的，我们直接返回所有数据，包含：
		  - 【响应头】200 OK
		  - 【响应头】Accept-Ranges : bytes
		  - 【响应头】Content-Type : video/mp4
		  - 【响应头】Content-Length : 1900
		  - 【响应头】ETag : "45b9ef0-5911062c10d4c"
		  - 【响应头】Last-Modified: Tue, 27 Aug 2019 02:54:05 GMT

		  - 【响应体】文件数据
		- 请求头中携带“Range”的，说明是进行断点续传了，我们返回的数据应当有所变化，包含：
		  - 【响应头】206 Partial Content
		  - 【响应头】Accept-Ranges : bytes
		  - 【响应头】Content-Type : video/mp4
		  - 【响应头】Content-Length : [需要传输的总长度]
		  - 【响应头】ETag : "45b9ef0-5911062c10d4c"
		  - 【响应头】Last-Modified : Tue, 27 Aug 2019 02:54:05 GMT
		  - 【响应头】Content-Range : bytes [start]-[end]/[total]

		  - 【响应体】文件数据
	*/

	// 通用响应头
	response_header := r.Response.Header()
	response_header["Accept-Ranges"] = append(response_header["Accept-Ranges"], "bytes")
	response_header["ETag"] = append(response_header["ETag"], etag)
	response_header["Last-Modified"] = append(response_header["Last-Modified"], filemodifytime.Format(time.RFC1123))

	// 文件传输块
	block_size := int64(4096 * 100)
	read_buf := make([]byte, block_size)

	real_file_object, err = OpenFile(real_path)
	if err != nil {
		r.Response.WriteStatus(500)
		return
	}

	if range_exist {
		// 断点续传，准备一下响应头
		Content_Range := fmt.Sprintf("bytes %d-%d/%d", start, end, filesize)
		response_header["Content-Range"] = append(response_header["Content-Range"], Content_Range)
		r.Response.WriteHeader(206)

		// 读取数据，发送出去
		_, err := real_file_object.Seek(start, 0)
		if err != nil {
			glog.Error(err)
			r.Response.WriteStatus(500)
			_ = real_file_object.Close()
			return
		}

		// 这里有可能导致申请失败
		read_buf = make([]byte, length)
		readed_bytes, err := real_file_object.Read(read_buf)
		if err != nil {
			glog.Error(err)
			r.Response.WriteStatus(500)
			_ = real_file_object.Close()
			return
		}

		// 这里执行解密

		r.Response.Write(read_buf[:readed_bytes])


	} else {
		// 读取数据，发送出去
		for {
			readed_bytes, err := real_file_object.Read(read_buf)
			if err != nil {
				//glog.Error(err)
				break
			}

			// 这里执行解密

			r.Response.Write(read_buf[:readed_bytes])
		}
	}

	_ = real_file_object.Close()
}