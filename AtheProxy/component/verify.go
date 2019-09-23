package component

import (
	"consts"
	"protocol/biz"
)

var Verify = &verify{}

type verify struct {
}

//检查参数
func (v *verify) CheckParams(serverBiz *serverbiz.ServerBiz) int32 {
	if serverBiz == nil {
		return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
	}
	//version
	if serverBiz.Version == "" {
		return int32(consts.MSP_ERROR_MSG_INVALID_VERSION)
	}
	//GlobalRoute
	if serverBiz.GlobalRoute == nil {
		return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
	}

	if serverBiz.GlobalRoute.SessionId == "" {
		return int32(consts.MSP_ERROR_MSG_SESSION_ID_EMPTY)
	}

	if serverBiz.GlobalRoute.Appid == "" {
		return int32(consts.MSP_ERROR_MSG_APP_ID_EMPTY)
	}

	switch int32(serverBiz.MsgType) {
	case int32(serverbiz.ServerBiz_UP_CALL):
		//UpCall
		if serverBiz.UpCall == nil {
			return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
		}

		if serverBiz.UpCall.Call == "" {
			return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
		}

		if serverBiz.UpCall.From == "" {
			return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
		}
		switch serverBiz.UpCall.From {
		case "":
			return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
		case "guider":
			if serverBiz.GlobalRoute.GuiderId == "" {
				return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
			}
		case "proxy":
			if serverBiz.GlobalRoute.UpRouterId == "" {
				return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
			}
		}
		return int32(consts.MSP_SUCCESS)
	case int32(serverbiz.ServerBiz_DOWN_RESULT):
		//TODO 暂不支持
		//DownResult
		return int32(consts.MSP_SUCCESS)
	default:
		return int32(consts.MSP_ERROR_MSG_PARAM_ERROR)
	}
}
