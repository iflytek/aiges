package consts

import (
	"errors"
)

const PREFIX = "/root/go/src/resource/config/"

//文件路径
const (
	XSFC_FILE   = "xsfc.toml"
	XSFS_FILE   = "xsfs.toml"
	CONFIG_FILE = "atmos.toml"
)

//错误消息
const (
	ERR_MSG_BAD_RPC         = "bad rpc call, data is nil"
	ERR_MSG_AUTH_NO_LICENSE = "licc failed"
	ERR_MSG_TIME_OUT        = "get result time out"
)

//错误码
const (
	MSP_SUCCESS               = 0
	MSP_ERROR_FAIL            = -1
	MSP_ERROR_EXCEPTION       = -2
	MSP_ATMOS_ERROR_EXCEPTION = -4

	/* General errors 10100(0x2774) */
	MSP_ERROR_GENERAL              = 10100 /* 0x2774 */
	MSP_ERROR_OUT_OF_MEMORY        = 10101 /* 0x2775 */
	MSP_ERROR_FILE_NOT_FOUND       = 10102 /* 0x2776 */
	MSP_ERROR_NOT_SUPPORT          = 10103 /* 0x2777 */
	MSP_ERROR_NOT_IMPLEMENT        = 10104 /* 0x2778 */
	MSP_ERROR_ACCESS               = 10105 /* 0x2779 */
	MSP_ERROR_INVALID_PARA         = 10106 /* 0x277A */
	MSP_ERROR_INVALID_PARA_VALUE   = 10107 /* 0x277B */
	MSP_ERROR_INVALID_HANDLE       = 10108 /* 0x277C */
	MSP_ERROR_INVALID_DATA         = 10109 /* 0x277D */
	MSP_ERROR_NO_LICENSE           = 10110 /* 0x277E */
	MSP_ERROR_NOT_INIT             = 10111 /* 0x277F */
	MSP_ERROR_NULL_HANDLE          = 10112 /* 0x2780 */
	MSP_ERROR_OVERFLOW             = 10113 /* 0x2781 */
	MSP_ERROR_TIME_OUT             = 10114 /* 0x2782 */
	MSP_ERROR_OPEN_FILE            = 10115 /* 0x2783 */
	MSP_ERROR_NOT_FOUND            = 10116 /* 0x2784 */
	MSP_ERROR_NO_ENOUGH_BUFFER     = 10117 /* 0x2785 */
	MSP_ERROR_NO_DATA              = 10118 /* 0x2786 */
	MSP_ERROR_NO_MORE_DATA         = 10119 /* 0x2787 */
	MSP_ERROR_NO_RESPONSE_DATA     = 10120 /* 0x2788 */
	MSP_ERROR_ALREADY_EXIST        = 10121 /* 0x2789 */
	MSP_ERROR_LOAD_MODULE          = 10122 /* 0x278A */
	MSP_ERROR_BUSY                 = 10123 /* 0x278B */
	MSP_ERROR_INVALID_CONFIG       = 10124 /* 0x278C */
	MSP_ERROR_VERSION_CHECK        = 10125 /* 0x278D */
	MSP_ERROR_CANCELED             = 10126 /* 0x278E */
	MSP_ERROR_INVALID_MEDIA_TYPE   = 10127 /* 0x278F */
	MSP_ERROR_CONFIG_INITIALIZE    = 10128 /* 0x2790 */
	MSP_ERROR_CREATE_HANDLE        = 10129 /* 0x2791 */
	MSP_ERROR_CODING_LIB_NOT_LOAD  = 10130 /* 0x2792 */
	MSP_ERROR_USER_CANCELLED       = 10131 /* 0x2793 */
	MSP_ERROR_INVALID_OPERATION    = 10132 /* 0x2794 */
	MSP_ERROR_MESSAGE_NOT_COMPLETE = 10133 /* 0x2795 */ //flash
	MSP_ERROR_NO_EID               = 10134 /* 0x2795 */
	MSP_ERROE_OVER_REQ             = 10135 /* 0x2797 */ //client Redundancy request
	MSP_ERROR_USER_ACTIVE_ABORT    = 10136 /* 0x2798 */ /*user abort*/
	MSP_ERROR_BUSY_GRMBUILDING     = 10137 /* 0x2799 */
	MSP_ERROR_BUSY_LEXUPDATING     = 10138 /* 0x279A */
	MSP_ERROR_SESSION_RESET        = 10139 //msc主动终止会话，准备重传

	/* Error codes of mssp message 10300(0x283C) */
	MSP_ERROR_MSG_GENERAL                = 10300 /* 0x283C */
	MSP_ERROR_MSG_PARSE_ERROR            = 10301 /* 0x283D */
	MSP_ERROR_MSG_BUILD_ERROR            = 10302 /* 0x283E */
	MSP_ERROR_MSG_PARAM_ERROR            = 10303 /* 0x283F */
	MSP_ERROR_MSG_CONTENT_EMPTY          = 10304 /* 0x2840 */
	MSP_ERROR_MSG_INVALID_CONTENT_TYPE   = 10305 /* 0x2841 */
	MSP_ERROR_MSG_INVALID_CONTENT_LENGTH = 10306 /* 0x2842 */
	MSP_ERROR_MSG_INVALID_CONTENT_ENCODE = 10307 /* 0x2843 */
	MSP_ERROR_MSG_INVALID_KEY            = 10308 /* 0x2844 */
	MSP_ERROR_MSG_KEY_EMPTY              = 10309 /* 0x2845 */
	MSP_ERROR_MSG_SESSION_ID_EMPTY       = 10310 /* 0x2846 */
	MSP_ERROR_MSG_LOGIN_ID_EMPTY         = 10311 /* 0x2847 */
	MSP_ERROR_MSG_SYNC_ID_EMPTY          = 10312 /* 0x2848 */
	MSP_ERROR_MSG_APP_ID_EMPTY           = 10313 /* 0x2849 */
	MSP_ERROR_MSG_EXTERN_ID_EMPTY        = 10314 /* 0x284A */
	MSP_ERROR_MSG_INVALID_CMD            = 10315 /* 0x284B */
	MSP_ERROR_MSG_INVALID_SUBJECT        = 10316 /* 0x284C */
	MSP_ERROR_MSG_INVALID_VERSION        = 10317 /* 0x284D */
	MSP_ERROR_MSG_NO_CMD                 = 10318 /* 0x284E */
	MSP_ERROR_MSG_NO_SUBJECT             = 10319 /* 0x284F */
	MSP_ERROR_MSG_NO_VERSION             = 10320 /* 0x2850 */
	MSP_ERROR_MSG_MSSP_EMPTY             = 10321 /* 0x2851 */
	MSP_ERROR_MSG_NEW_RESPONSE           = 10322 /* 0x2852 */
	MSP_ERROR_MSG_NEW_CONTENT            = 10323 /* 0x2853 */
	MSP_ERROR_MSG_INVALID_SESSION_ID     = 10324 /* 0x2854 */
	MSP_ERROR_MSG_INVALID_CONTENT        = 10325 /* 0x2855 */

	/* Error codes of DataBase 10400(0x28A0)*/
	MSP_ERROR_DB_GENERAL       = 10400 /* 0x28A0 */
	MSP_ERROR_DB_EXCEPTION     = 10401 /* 0x28A1 */
	MSP_ERROR_DB_NO_RESULT     = 10402 /* 0x28A2 */
	MSP_ERROR_DB_INVALID_USER  = 10403 /* 0x28A3 */
	MSP_ERROR_DB_INVALID_PWD   = 10404 /* 0x28A4 */
	MSP_ERROR_DB_CONNECT       = 10405 /* 0x28A5 */
	MSP_ERROR_DB_INVALID_SQL   = 10406 /* 0x28A6 */
	MSP_ERROR_DB_INVALID_APPID = 10407 /* 0x28A7 */
	MSP_ERROR_DB_NO_UID        = 10408

	/* Error Codes using in local engine */
	MSP_ERROR_AUTH_NO_LICENSE            = 11200 /* 0x2BC0 */ /* 无授权 */
	MSP_ERROR_AUTH_NO_ENOUGH_LICENSE     = 11201 /* 0x2BC1 */ /* 授权不足 */
	MSP_ERROR_AUTH_INVALID_LICENSE       = 11202 /* 0x2BC2 */ /* 无效的授权 */
	MSP_ERROR_AUTH_LICENSE_EXPIRED       = 11203 /* 0x2BC3 */ /* 授权过期 */
	MSP_ERROR_AUTH_NEED_MORE_DATA        = 11204 /* 0x2BC4 */ /* 无设备信息 */
	MSP_ERROR_AUTH_LICENSE_TO_BE_EXPIRED = 11205 /* 0x2BC5 */ /* 授权即将过期，警告性错误码 */

	/* Error codes of http 12000(0x2EE0) */
	MSP_ERROR_HTTP_BASE = 12000 /* 0x2EE0 */
	MSP_ERROR_HTTP_400  = 12400
	MSP_ERROR_HTTP_401  = 12401
	MSP_ERROR_HTTP_402  = 12402
	MSP_ERROR_HTTP_403  = 12403
	MSP_ERROR_HTTP_404  = 12404
	MSP_ERROR_HTTP_405  = 12405
	MSP_ERROR_HTTP_406  = 12406
	MSP_ERROR_HTTP_407  = 12407
	MSP_ERROR_HTTP_408  = 12408
	MSP_ERROR_HTTP_409  = 12409
	MSP_ERROR_HTTP_410  = 12410
	MSP_ERROR_HTTP_411  = 12411
	MSP_ERROR_HTTP_412  = 12412
	MSP_ERROR_HTTP_413  = 12413
	MSP_ERROR_HTTP_414  = 12414
	MSP_ERROR_HTTP_415  = 12415
	MSP_ERROR_HTTP_416  = 12416
	MSP_ERROR_HTTP_417  = 12417
	MSP_ERROR_HTTP_500  = 12500
	MSP_ERROR_HTTP_501  = 12501
	MSP_ERROR_HTTP_502  = 12502
	MSP_ERROR_HTTP_503  = 12503
	MSP_ERROR_HTTP_504  = 12504
	MSP_ERROR_HTTP_505  = 12505
)

const SERVER_ERROR_TIME_OUT = 60114 /*  调用引擎获取结果超时 */

//引擎常量
const (
	SESSION_HANDLE = "ent_session_handle"
	ENT_OP_IN      = "AIIn"
	ENT_OP_OUT     = "AIOut"
	ENT_OP_EXP     = "AIExcp"
)

const SERVICE_NAME = "serviceName"

const SESSIONID = "sid"

//callback 状态
const TASK_STATUS = "task_status"

//callback 错误描述信息
const ERROR_INFO = "error_info"

//callback 成功状态
const STATUS_SUCCEED = "SUCCEED"

//callback 处理中状态
const STATUS_IN_PROCESS = "IN_PROCESS"

//callback 失败状态
const STATUS_FAILED = "FAILED"

//callback 失败状态
const ATMOS_ROUTE = "engine_route"

const TIME_OUT = "time out"

//调用来源
const SOURCE_GUIDER = "guider"

//0表示关，1表示开
const OPEN = 1

const DEMOTE_IS_OPEN = "1"

//0表示关，1表示开
const CLOSE = 0

const LOG_LEVEL_INFO = "INFO"

const LOG_LEVEL_DEBUG = "DEBUG"

//个性化NLP result
const NLP_RESULT = "nlp_result"

const DEFAULT_TIMEOUT = "defaultTimeout"

const ENGINE_TIMEOUT = "engineTimeout"

//个性化开关，rp=1 个性化分段
const RP = "rp"

//可以忽略的错误码：这类错误码，不需要通知引擎，会话出错
var IGNORE_ERROR_ARRAY = []int32{10101}

const OP_EXCEPTION = "exception"

var NO_LICENSE_ERR = errors.New(ERR_MSG_AUTH_NO_LICENSE)

var TIME_OUT_ERR = errors.New(ERR_MSG_TIME_OUT)

var MSG_BAD_RPC = errors.New(ERR_MSG_BAD_RPC)

var ERR = errors.New("")

var ERR_PARAM_INVALID = errors.New("param invalid")

const AIAAS_RENAME_PRX = "aiaas_rename"
const LOCAL_IP = "127.0.0.1"
const DEFAULT_PORT = "8686"
const DEFAULT_EVENT_PORT = "4545"
