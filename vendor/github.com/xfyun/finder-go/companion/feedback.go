package finder

import (
	"encoding/json"
	"fmt"
	"net/http"

	common "github.com/xfyun/finder-go/common"
	errors "github.com/xfyun/finder-go/errors"
	"github.com/xfyun/finder-go/utils/httputil"

	"github.com/xfyun/finder-go/log"
)

func FeedbackForConfig(hc *http.Client, url string, f *common.ConfigFeedback) error {
	log.Log.Infof("call FeedbackForConfig")
	contentType := "application/x-www-form-urlencoded"
	params := []byte(fmt.Sprintf("push_id=%s&project=%s&group=%s&service=%s&version=%s&addr=%s&config=%s&update_time=%d&update_status=%d&load_time=%d&load_status=%d&gray_group_id=%s",
		f.PushID, f.ServiceMete.Project, f.ServiceMete.Group, f.ServiceMete.Service, f.ServiceMete.Version, f.ServiceMete.Address,
		f.Config, f.UpdateTime, f.UpdateStatus, f.LoadTime, f.LoadStatus, f.GrayGroupId))
	result, err := httputil.DoPost(hc, contentType, url, params)
	if err != nil {
		log.Log.Errorf("FeedbackForConfig err: %s", err)
		err = errors.NewFinderError(errors.FeedbackPostErr)
		return err
	} else {
		log.Log.Infof("FeedbackForConfig result: %s", string(result))
	}

	var r JSONResult
	err = json.Unmarshal([]byte(result), &r)
	if err != nil {
		log.Log.Errorf("FeedbackForConfig err: %s", err)
		err = errors.NewFinderError(errors.JsonUnmarshalErr)
		return err
	}
	if r.Ret != 0 {
		err = errors.NewFinderError(errors.FeedbackConfigErr)
		log.Log.Errorf("FeedbackForConfig err: %s", r.Msg)
		return err
	}

	return nil
}

func FeedbackForService(hc *http.Client, url string, f *common.ServiceFeedback) error {
	contentType := "application/x-www-form-urlencoded"
	params := []byte(fmt.Sprintf("push_id=%s&project=%s&group=%s&consumer=%s&consumer_version=%s&addr=%s&provider=%s&provider_version=%s&update_time=%d&update_status=%d&load_time=%d&load_status=%d&api_version=%s&type=%d",
		f.PushID, f.ServiceMete.Project, f.ServiceMete.Group, f.ServiceMete.Address, f.ServiceMete.Version, f.ServiceMete.Address,
		f.Provider, f.ProviderVersion, f.UpdateTime, f.UpdateStatus, f.LoadTime, f.LoadStatus, f.ProviderVersion, f.Type))
	result, err := httputil.DoPost(hc, contentType, url, params)
	if err != nil {
		log.Log.Errorf("%s", err)
		err = errors.NewFinderError(errors.FeedbackPostErr)
		return err
	} else {
		log.Log.Infof("FeedbackForService result: %s", string(result))
	}

	var r JSONResult
	err = json.Unmarshal([]byte(result), &r)
	if err != nil {
		log.Log.Errorf("[FeedbackForService][json] %s", err)
		err = errors.NewFinderError(errors.JsonUnmarshalErr)
		return err
	}
	if r.Ret != 0 {
		err = errors.NewFinderError(errors.FeedbackServiceErr)
		log.Log.Errorf("FeedbackServiceError : %s", r.Msg)
		return err
	}

	return nil
}
