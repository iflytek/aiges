package request

import (
	"fmt"
	"git.iflytek.com/AIaaS/xsf/utils"
	"io/ioutil"
	"os"
	"sync"
)

// 下行数据异步落盘或打印
func (r *Request) DownStreamWrite(wg *sync.WaitGroup, log *utils.Logger) {

	for {
		output, alive := <-r.C.AsyncDrop
		if !alive {
			break // channel 关闭, 退出
		}

		key := output.Sid + "-" + output.Type + "-" + output.Name
		r.downOutput(key, output.Data, log)
	}
	wg.Done()
}

func (r *Request) downOutput(key string, data []byte, log *utils.Logger) {
	switch r.C.Output {
	case 0:
		fi, err := os.OpenFile(r.C.OutputDst, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Errorw("downOutput Sync OpenFile fail", "err", err.Error(), "key", key)
			return
		}

		tmp := []byte(key + ":")
		tmp = append(tmp, data...)
		tmp = append(tmp, byte('\n'))
		wlen, err := fi.Write(tmp)
		if err != nil || wlen != len(tmp) {
			log.Errorw("downOutput Sync AppendFile fail", "err", err.Error(), "wlen", wlen, "key", key)
			_ = fi.Close()
			return
		}
		_ = fi.Close()
	case 1: // 输出至目录OutputDst
		err := ioutil.WriteFile(r.C.OutputDst+"/"+key, data, os.ModePerm)
		if err != nil {
			log.Errorw("downOutput Sync WriteFile fail", "err", err.Error(), "key", key)
			return
		}
	case 2: // 输出至终端
		fmt.Println(key, ":", string(data))
	case -1:
		// 下行数据不输出, nothing to do
	}
}
