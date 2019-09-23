package storage

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// 封装http协议,用于与外部业务数据交互的上传下载;
func HttpDownload(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		return data, errors.New("http download fail with " + resp.Status)
	}
	data, err = ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return
}
