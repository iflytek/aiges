package common

type ErrorResp struct {
	Header Header `json:"header"`
}

type Header struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Sid     string    `json:"sid"`
	Status  int       `json:"status"`
	Cid     string    `json:"cid,omitempty"`
}

type SuccessResp struct {
	//Header Header `json:"header"`
	Header   map[string]interface{} `json:"header"`
	Payload  interface{}            `json:"payload,omitempty"`
	WfStatus int                    `json:"wf_status,omitempty"`
}

func (s *SuccessResp) SetHeader(k string, v interface{}) {
	if s.Header == nil {
		s.Header = map[string]interface{}{}
	}
	if !isEmpty(v) {
		s.Header[k] = v
	}
}

func isEmpty(v interface{}) bool {
	switch v := v.(type) {
	case string:
		return v == ""
	case nil:
		return true
	}
	return false
}

type HttpSuccessResp struct {
	//Code int `json:"code"`
	//Message string `json:"message"`
	//Sid string `json:"sid"`
	Header  map[string]interface{} `json:"header"`
	Payload interface{}            `json:"payload,omitempty"`
}

func NewSuccessResp(sid string, payload interface{}, status int, cid string) *SuccessResp {
	resp := &SuccessResp{
		Header: map[string]interface{}{
			"code":    0,
			"message": "success",
			"sid":     sid,
			"status":  status,
		},
		//Header: Header{
		//	Code:    0,
		//	Message: "success",
		//	Sid:     sid,
		//	Status:  status,
		//	Cid: cid,
		//},
		Payload: payload,
	}
	resp.SetHeader("cid", cid)
	return resp
}

//	func NewHttpSuccessResp(sid string, payload interface{}, s *HttpSession) *HttpSuccessResp {
//		hd := map[string]interface{}{}
//		for key, val := range s.schema.BuildResponseHeader(s.respHeader) {
//			hd[key] = val
//		}
//		hd["code"] = 0
//		hd["message"] = "success"
//		hd["sid"] = sid
//		if s.schema.Meta.EnableClientSession() {
//			hd["session"] = s.ClientSession
//		}
//		return &HttpSuccessResp{
//			Header: hd,
//			//Header:Header{
//			//	Code:    0,
//			//	Message: "success",
//			//	Sid:     sid,
//			//},
//			Payload: payload,
//		}
//	}
func NewHttpErrorResp(sid string, code int, msg string) *HttpSuccessResp {
	return &HttpSuccessResp{
		Header: map[string]interface{}{
			"code":    code,
			"message": msg,
			"sid":     sid,
		},
	}
}
