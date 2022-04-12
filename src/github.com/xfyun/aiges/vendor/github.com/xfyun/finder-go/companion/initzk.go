package finder

import (
	"encoding/json"
	"net/http"
	"time"

	common "github.com/xfyun/finder-go/common"
	errors "github.com/xfyun/finder-go/errors"
	"github.com/xfyun/finder-go/utils/httputil"
)

// GetStorageInfo for getting storage metadata
func GetStorageInfo(hc *http.Client, url string) (*common.StorageInfo, error) {
	var result []byte
	var err error
	retryNum := 0
	for {
		result, err = httputil.DoGet(hc, url)
		//log.Println("向")
		if err != nil {
			if retryNum < 3 {
				retryNum++
				time.Sleep(time.Millisecond * 100)
				continue
			} else {
				return nil, err
			}
		} else {
			break
		}
	}

	var r JSONResult
	err = json.Unmarshal([]byte(result), &r)
	if err != nil {
		return nil, err
	}
	if r.Ret != 0 {
		err = errors.NewFinderError(errors.ZkGetInfoError)
		return nil, err
	}

	ok := true
	if _, ok = r.Data["config_path"]; !ok {
		err = errors.NewFinderError(errors.ZkInfoMissConfigRootPath)
		return nil, err
	}

	if _, ok = r.Data["service_path"]; !ok {
		err = errors.NewFinderError(errors.ZkInfoMissServiceRootPath)
		return nil, err
	}
	if _, ok := r.Data["zk_node_path"]; !ok {
		err = errors.NewFinderError(errors.ZkInfoMissZkNodePath)
		return nil, err
	}
	if r.Data["zk_node_path"] == nil {
		err = errors.NewFinderError(errors.ZkInfoMissZkNodePath)
		return nil, err
	}
	var zkAddr []string
	if _, ok = r.Data["zk_addr"]; !ok {
		err = errors.NewFinderError(errors.ZkInfoMissAddr)
		return nil, err
	}

	var value []interface{}
	if value, ok = r.Data["zk_addr"].([]interface{}); ok {
		zkAddr = convert(value)
		if len(zkAddr) == 0 {
			err = errors.NewFinderError(errors.ZkInfoAddrConvertError)
			return nil, err
		}
	}

	zkInfo := &common.StorageInfo{
		ConfigRootPath:  r.Data["config_path"].(string),
		ServiceRootPath: r.Data["service_path"].(string),
		Addr:            zkAddr,
		ZkNodePath:      r.Data["zk_node_path"].(string),
	}

	return zkInfo, nil
}

func convert(in []interface{}) []string {
	r := make([]string, 0)
	ok := true
	value := ""
	for _, v := range in {
		if value, ok = v.(string); ok {
			r = append(r, value)
		}
	}

	return r
}
