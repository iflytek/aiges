package frame

const (
	ParaSubject	  = "sub"		// 请求服务类别
	ParaAudioRate = "rate"		// 音频数据采样率
	ParaAppId	  = "appid"		// 应用id
	ParaUsrId	  = "uid"		// 用户id
	ParaSessId	  = "sid"		// 会话session id
	ParaAuCodec	  = "encoding"	// 音频编解码类型
	ParaSpxFrame  = "spx_fsize"	// 开源speex编码帧大小(适配老参数,风格未统一)
	ParaNrtTask	  = "task_id"	// 非实时任务,用于更新任务状态(适配上层参数,风格未统一)
)
