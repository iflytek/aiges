package codec

import (
	"github.com/xfyun/aiges/frame"
	"strconv"
)

// 实现视频编解码功能
var vdcodecs map[string]codecor

func init() {
	vdcodecs = make(map[string]codecor)
	ffmpegInit()
}

type vdcdesc struct {
	enc       string // 编解码类型
	frameRate int    // frame_rate
	width     int
	heigth    int
}

type VdcodecInst struct {
	inst interface{} // 编解码实例
	vdcdesc
}

func (vi *VdcodecInst) Encode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := vdcodecs[vi.enc]; codec != nil {
		if output, code, err = codec.encode(vi.inst, input, last); err != nil {
			code = frame.AigesErrorCodecEncode
		}
		return
	}
	output = input
	return
}

func (vi *VdcodecInst) Decode(input []byte, last bool) (output []byte, code int, err error) {
	/*
	if codec := vdcodecs[vi.enc]; codec != nil {
		if output, code, err = codec.decode(vi.inst, input, last); err != nil {
			code = frame.AigesErrorCodecDecode
		}
		return
	}
	*/
	output = input
	return
}

func NewVdcodec(tag VideoTag) (vi *VdcodecInst, code int, err error) {
	vi = &VdcodecInst{vdcdesc: vdcdesc{enc: tag.Encoding}}
	vi.heigth, _ = strconv.Atoi(tag.Height)
	vi.width, _ = strconv.Atoi(tag.Width)
	vi.frameRate, _ = strconv.Atoi(tag.FrameRate)
	if codec := vdcodecs[vi.enc]; codec != nil {
		if vi.inst, code, err = codec.start(&vi.vdcdesc); err != nil {
			code = frame.AigesErrorCodecStart
		}
	}
	return
}

func CloseVdcodec(vi *VdcodecInst) (code int, err error) {
	if vi != nil {
		if codec := vdcodecs[vi.enc]; codec != nil {
			if code, err = codec.stop(vi.inst); err != nil {
				code = frame.AigesErrorCodecStop
			}
		}
	}
	return
}
