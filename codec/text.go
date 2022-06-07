package codec

import "github.com/xfyun/aiges/frame"

// 实现文本编码功能
var txtcodecs map[string]codecor

func init() {
	txtcodecs = make(map[string]codecor)
}

type txtdesc struct {
	encoding string // 编码类型
	compress string // 压缩类型
}

type TxtcodecInst struct {
	inst interface{} // 编解码实例
	txtdesc
}

func (ti *TxtcodecInst) Encode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := txtcodecs[ti.encoding]; codec != nil {
		if output, code, err = codec.encode(ti.inst, input, last); err != nil {
			code = frame.AigesErrorCodecEncode
		}
		return
	}
	output = input
	return
}

func (ti *TxtcodecInst) Decode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := txtcodecs[ti.encoding]; codec != nil {
		if output, code, err = codec.decode(ti.inst, input, last); err != nil {
			code = frame.AigesErrorCodecDecode
		}
		return
	}
	output = input
	return
}

func NewTxtcodec(tag TextTag) (ti *TxtcodecInst, code int, err error) {
	ti = &TxtcodecInst{txtdesc: txtdesc{tag.Encoding, tag.Compress}}
	if codec := txtcodecs[ti.encoding]; codec != nil {
		ti.inst, code, err = codec.start(&ti.txtdesc)
	}
	return
}

func CloseTxtcodec(ti *TxtcodecInst) (code int, err error) {
	if ti != nil {
		if codec := txtcodecs[ti.encoding]; codec != nil {
			code, err = codec.stop(ti.inst)
		}
	}
	return
}
