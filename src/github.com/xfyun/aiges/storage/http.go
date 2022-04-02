package storage

import (
 	"errors"
	"github.com/xfyun/aiges/conf"
	"github.com/xfyun/aiges/frame"
 	aigesUtils "github.com/xfyun/aiges/utils"
	"io/ioutil"
	"net/http"
)

// 封装http协议,用于与外部业务数据交互的上传下载;
func HttpDownload(url string) (data []byte, code int, err error) {
	defer func() {
		aigesUtils.CommonLogger.Debugw("resource http download ",
			"url", url, "dataLength", len(data), "code", code, "err ", err)
	}()

	endFlag := make(chan bool)

	for i := 0; i < conf.HttpRetry; i++ {
		go func() {
			defer func() { endFlag <- true }()
			var resp *http.Response
			resp, err = http.Get(url)
			if err != nil {
				code = frame.AigesErrorHttpReq
				err = errors.New("http download fail with" + err.Error())
				return
			} else if resp.StatusCode != 200 {
				code = frame.AigesErrorHttpFail
				err = errors.New("http download fail with " + resp.Status)
				_ = resp.Body.Close()
				return
			}
			data, err = ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				code = frame.AigesErrorHttpInvalidData
				err = errors.New("http download fail with" + err.Error())
				return
			}
		}()
		select {
		case <-endFlag:
			if code != 0 {
				continue
			}
			return
		}
	}
	return
}
