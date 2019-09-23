package finder

import (
	"encoding/json"
	"log"

	"fmt"
	"net/http"

	errors "git.xfyun.cn/AIaaS/finder-go/errors"
	"git.xfyun.cn/AIaaS/finder-go/utils/httputil"
)

func RegisterService(hc *http.Client, url string, project string, group string, service string, apiVersion string) error {
	contentType := "application/x-www-form-urlencoded"
	params := []byte(fmt.Sprintf("project=%s&group=%s&service=%s&api_version=%s", project, group, service, apiVersion))
	result, err := httputil.DoPost(hc, contentType, url, params)
	if err != nil {
		log.Println(err)
		return err
	}

	var r JSONResult
	err = json.Unmarshal([]byte(result), &r)
	if err != nil {
		return err
	}
	if r.Ret != 0 {
		log.Println("向companion注册服务失败：companion返回失败：",r.Msg," code:",r.Ret)
		err = errors.NewFinderError(errors.CompanionRegisterServiceErr)
		return err
	}

	return nil
}
