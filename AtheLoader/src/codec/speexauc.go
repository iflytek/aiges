package codec

/*
#cgo CFLAGS: -I../cgo/header
#cgo LDFLAGS: -L../cgo/library -lspeex1.5.1
#include <stdlib.h>
#include "./speex/speex.h"
*/
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

// 开源版本speex编解码类型,带缓存功能的speex解码器
type speexInst struct {
	bits           C.SpeexBits
	decState       unsafe.Pointer
	frameSize      int
	spxDecodeFrame int    // default encode quality is 8 (ranging from 0 to 10 )
	cacheBuf       []byte // cache buf to storage loss frame

}

type speexCodec struct {
}

func (sc *speexCodec) start(desc *aucdesc) (inst interface{}, code int, err error) {
	speex := &speexInst{}
	C.speex_bits_init(&speex.bits)
	if desc.auc == AUDIOSPEEXRAWWB {
		speex.decState = C.speex_decoder_init(&C.speex_wb_mode)
	} else {
		speex.decState = C.speex_decoder_init(&C.speex_nb_mode)
	}
	// codec option control
	enh := C.int(1)
	sr, _ := strconv.Atoi(desc.sampleRate)
	rate := C.int(sr)
	var frame C.int
	C.speex_decoder_ctl(speex.decState, C.SPEEX_SET_ENH, unsafe.Pointer(&enh))
	C.speex_decoder_ctl(speex.decState, C.SPEEX_SET_SAMPLING_RATE, unsafe.Pointer(&rate))
	C.speex_encoder_ctl(speex.decState, C.SPEEX_GET_FRAME_SIZE, unsafe.Pointer(&frame))
	speex.frameSize = int(frame)
	speex.spxDecodeFrame, _ = strconv.Atoi(desc.frame)
	// buffer initial
	if speex.spxDecodeFrame > 0 {
		speex.cacheBuf = make([]byte, 0, speex.spxDecodeFrame)
	}
	inst = speex
	return
}

func (sc *speexCodec) stop(inst interface{}) (code int, err error) {
	speex := inst.(*speexInst)
	C.speex_bits_destroy(&(speex.bits))
	if speex.decState != nil {
		C.speex_decoder_destroy(speex.decState)
	}
	return
}

func (sc *speexCodec) encode(inst interface{}, input []byte, last bool) (output []byte, code int, err error) {
	return nil, -1, errors.New("not supported yet")
}

func (sc *speexCodec) decode(inst interface{}, input []byte, last bool) (output []byte, code int, err error) {
	speex := inst.(*speexInst)
	// 固定帧解码
	if speex.spxDecodeFrame > 0 {
		readLen := 0
		cacheSpare := speex.spxDecodeFrame - len(speex.cacheBuf)
		for readLen < len(input) {
			cpLen := cacheSpare
			if cacheSpare >= len(input)-readLen {
				cpLen = len(input) - readLen
			}
			speex.cacheBuf = append(speex.cacheBuf, input[readLen:readLen+cpLen]...)
			readLen += cpLen
			// 存在整帧或已至尾帧,解码
			if len(speex.cacheBuf) == speex.spxDecodeFrame || last {
				audioBuf := make([]byte, C.sizeof_spx_int16_t*speex.frameSize, C.sizeof_spx_int16_t*speex.frameSize)
				C.speex_bits_read_from(&speex.bits, (*C.char)(unsafe.Pointer(&speex.cacheBuf[0])), C.int(len(speex.cacheBuf)))
				C.speex_decode_int(speex.decState, &speex.bits, (*C.spx_int16_t)(unsafe.Pointer(&audioBuf[0])))
				output = append(output, audioBuf...)
				// 清空缓存
				speex.cacheBuf = speex.cacheBuf[0:0]
				cacheSpare = speex.spxDecodeFrame
			}
		}
	} else {
		// TODO 非固定帧解码
	}

	return
}
