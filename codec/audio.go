package codec

import "github.com/xfyun/aiges/frame"

var aucodecs map[string]codecor

func init() {
	aucodecs = make(map[string]codecor)
}

type aucdesc struct {
	auc          string // 编解码类型
	compressRate string // 压缩率,不同编解码类型(如speex:0~10)
	sampleRate   string // 采样率,如8000/16000等
	params       string // 编解码参数
	frame        string // 编解码压缩帧大小
}

type AucodecInst struct {
	inst interface{} // 编解码实例
	aucdesc
}

func (ai *AucodecInst) Encode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := aucodecs[ai.auc]; codec != nil {
		if output, code, err = codec.encode(ai.inst, input, last); err != nil {
			code = frame.AigesErrorCodecEncode
		}
		return
	}
	output = input
	return
}

func (ai *AucodecInst) Decode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := aucodecs[ai.auc]; codec != nil {
		if output, code, err = codec.decode(ai.inst, input, last); err != nil {
			code = frame.AigesErrorCodecDecode
		}
		return
	}
	output = input
	return
}

func NewAucodec(tag AudioTag) (ai *AucodecInst, code int, err error) {
	var desc aucdesc
	if desc, code, err = normalize(tag); err == nil {
		ai = &AucodecInst{aucdesc: desc}
		if codec := aucodecs[ai.auc]; codec != nil {
			ai.inst, code, err = codec.start(&ai.aucdesc)
		}
	}
	return
}

func CloseAucodec(ai *AucodecInst) (code int, err error) {
	if ai != nil {
		if codec := aucodecs[ai.auc]; codec != nil {
			code, err = codec.stop(ai.inst)
		}
	}
	return
}
