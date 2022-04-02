package codec

import "github.com/xfyun/aiges/frame"

// 实现图像编解码功能
var imgcodecs map[string]codecor

func init() {
	imgcodecs = make(map[string]codecor)
}

type imgdesc struct {
	encoding string // 编码类型
}

type ImgcodecInst struct {
	inst interface{} // 编解码实例
	imgdesc
}

func (ti *ImgcodecInst) Encode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := imgcodecs[ti.encoding]; codec != nil {
		if output, code, err = codec.encode(ti.inst, input, last); err != nil {
			code = frame.AigesErrorCodecEncode
		}
		return
	}
	output = input
	return
}

func (ti *ImgcodecInst) Decode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := imgcodecs[ti.encoding]; codec != nil {
		if output, code, err = codec.decode(ti.inst, input, last); err != nil {
			code = frame.AigesErrorCodecDecode
		}
		return
	}
	output = input
	return
}

func NewImgcodec(tag ImageTag) (ti *ImgcodecInst, code int, err error) {
	ti = &ImgcodecInst{imgdesc: imgdesc{tag.Encoding}}
	if codec := imgcodecs[ti.encoding]; codec != nil {
		ti.inst, code, err = codec.start(&ti.imgdesc)
	}
	return
}

func CloseImgcodec(ti *ImgcodecInst) (code int, err error) {
	if ti != nil {
		if codec := imgcodecs[ti.encoding]; codec != nil {
			code, err = codec.stop(ti.inst)
		}
	}
	return
}
