package codec

/*
#cgo CFLAGS: -I../cgo/header
#cgo LDFLAGS: -L../cgo/library -lspeex1.5.1
#include <stdlib.h>
#include "./speex/speex.h"
*/
import "C"
import (
	"strconv"
	"unsafe"
)

// 开源版本speex编解码类型,带缓存功能的speex解码器
type speexInst struct {
	bits           C.SpeexBits
	decState       unsafe.Pointer
	encState       unsafe.Pointer
	frameSize      int
	spxDecodeFrame int    // default encode quality is 8 (ranging from 0 to 10 )
	decCacheBuf    []byte // cache buf to storage loss frame
	encCacheBuf    []byte // cache buf to storage loss frame

	from int //
	to   int
}

type speexCodec struct {
}

func (sc *speexCodec) start(it interface{}) (inst interface{}, code int, err error) {
	desc := it.(*aucdesc)
	speex := &speexInst{}
	C.speex_bits_init(&speex.bits)

	if desc.auc == AUDIOSPEEXRAWWB {
		// decode
		speex.decState = C.speex_decoder_init(&C.speex_wb_mode)
		// encode
		rate := C.int(16000)
		speex.encState = C.speex_encoder_init(&C.speex_wb_mode)
		C.speex_encoder_ctl(speex.encState, C.SPEEX_SET_SAMPLING_RATE, unsafe.Pointer(&rate))
	} else {
		// decode
		speex.decState = C.speex_decoder_init(&C.speex_nb_mode)
		// encode
		rate := C.int(8000)
		speex.encState = C.speex_encoder_init(&C.speex_nb_mode)
		C.speex_encoder_ctl(speex.encState, C.SPEEX_SET_SAMPLING_RATE, unsafe.Pointer(&rate))
	}

	// decode option control
	enh := C.int(1)
	sr, _ := strconv.Atoi(desc.sampleRate)
	rate := C.int(sr)
	var frame C.int
	C.speex_decoder_ctl(speex.decState, C.SPEEX_SET_ENH, unsafe.Pointer(&enh))
	C.speex_decoder_ctl(speex.decState, C.SPEEX_SET_SAMPLING_RATE, unsafe.Pointer(&rate))

	// get SPEEX_GET_FRAME_SIZE
	C.speex_encoder_ctl(speex.decState, C.SPEEX_GET_FRAME_SIZE, unsafe.Pointer(&frame))
	speex.frameSize = int(frame)

	// decode buffer initial
	speex.spxDecodeFrame, _ = strconv.Atoi(desc.frame)
	if speex.spxDecodeFrame > 0 {
		speex.decCacheBuf = make([]byte, 0, speex.spxDecodeFrame)
	}

	// encode
	r, e := strconv.Atoi(desc.compressRate)
	if e != nil {
		r = 8
	}
	qua := C.int(r)
	C.speex_encoder_ctl(speex.encState, C.SPEEX_SET_QUALITY, unsafe.Pointer(&qua))
	speex.encCacheBuf = make([]byte, 0, C.sizeof_spx_int16_t*speex.frameSize)
	C.speex_bits_init(&speex.bits)

	inst = speex
	return
}

func (sc *speexCodec) stop(inst interface{}) (code int, err error) {
	speex := inst.(*speexInst)
	C.speex_bits_destroy(&(speex.bits))
	if speex.decState != nil {
		C.speex_decoder_destroy(speex.decState)
		C.speex_encoder_destroy(speex.encState)
	}
	return
}

func (sc *speexCodec) encode(inst interface{}, input []byte, last bool) (output []byte, code int, err error) {
	speex := inst.(*speexInst)
	frm := C.sizeof_spx_int16_t * speex.frameSize
	// fmt.Println(speex.frameSize, frm)

	var inputIndex int
	cacheSpare := frm - len(speex.encCacheBuf) //剩余空间

	for inputIndex < len(input) {
		if cacheSpare <= (len(input) - inputIndex) {
			speex.encCacheBuf = append(speex.encCacheBuf, input[inputIndex:inputIndex+cacheSpare]...)
			inputIndex += cacheSpare

			C.speex_bits_reset(&speex.bits)
			C.speex_encode_int(speex.encState, (*C.spx_int16_t)(unsafe.Pointer(&speex.encCacheBuf[0])), &speex.bits)
			out := make([]byte, frm, frm)
			nb := C.speex_bits_write(&speex.bits, (*C.char)(unsafe.Pointer(&out[0])), (C.int)(frm))

			speex.from = frm
			speex.to = int(nb)
			output = append(output, out[:nb]...)
			// fmt.Println("enc", len(speex.encCacheBuf), "->", nb)

			speex.encCacheBuf = speex.encCacheBuf[0:0]
			cacheSpare = frm
		} else {
			speex.encCacheBuf = input[inputIndex:]
			// fmt.Println("cache", len(input)-inputIndex)
			speex.from = len(input) - inputIndex
			speex.to = 0
			break
		}
	}
	return

	//return nil, -1, errors.New("not supported yet")

}

func (sc *speexCodec) decode(inst interface{}, input []byte, last bool) (output []byte, code int, err error) {
	speex := inst.(*speexInst)
	// 固定帧解码
	if speex.spxDecodeFrame > 0 {
		readLen := 0
		cacheSpare := speex.spxDecodeFrame - len(speex.decCacheBuf)
		for readLen < len(input) {
			cpLen := cacheSpare
			if cacheSpare >= len(input)-readLen {
				cpLen = len(input) - readLen
			}
			speex.decCacheBuf = append(speex.decCacheBuf, input[readLen:readLen+cpLen]...)
			readLen += cpLen
			// 存在整帧或已至尾帧,解码
			if len(speex.decCacheBuf) == speex.spxDecodeFrame || last {
				audioBuf := make([]byte, C.sizeof_spx_int16_t*speex.frameSize, C.sizeof_spx_int16_t*speex.frameSize)
				C.speex_bits_read_from(&speex.bits, (*C.char)(unsafe.Pointer(&speex.decCacheBuf[0])), C.int(len(speex.decCacheBuf)))
				C.speex_decode_int(speex.decState, &speex.bits, (*C.spx_int16_t)(unsafe.Pointer(&audioBuf[0])))
				output = append(output, audioBuf...)
				// 清空缓存
				speex.decCacheBuf = speex.decCacheBuf[0:0]
				cacheSpare = speex.spxDecodeFrame
			}
		}
	} else {
		// TODO 非固定帧解码
	}

	return
}
