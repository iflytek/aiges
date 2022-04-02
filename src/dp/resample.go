package dp

/*
#cgo CFLAGS: -I../cgo/header
#cgo LDFLAGS: -L../cgo/library -lspeexdsp
#include <stdlib.h>
#include <stdint.h>
#include "./speexdsp/speex_resampler.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/xfyun/aiges/frame"
	"unsafe"
)

type AudioResampler struct {
	instance *C.SpeexResamplerState
	inRate   int
	outRate  int
}

func (rs *AudioResampler) Init(channels int, inRate int, outRate int, quality int) error {
	var errNum C.int
	rs.instance = C.speex_resampler_init(C.spx_uint32_t(channels), C.spx_uint32_t(inRate),
		C.spx_uint32_t(outRate), C.int(quality), &errNum)
	if rs.instance == nil {
		errInfo := fmt.Sprintf(" err: %d, chan:%d, inRate:%d, outRate:%d, quality:%d",
			int(errNum), channels, inRate, outRate, quality)
		return errors.New(frame.ErrorAudioResampleInit.Error() + errInfo)
	}
	rs.inRate = inRate
	rs.outRate = outRate
	return nil
}

func (rs *AudioResampler) ProcessInt(chanIndex int, bufIn []byte) (bufOut []byte, err error) {
	if rs.instance != nil && len(bufIn) != 0 {
		bufOut = make([]byte, len(bufIn)*rs.outRate/rs.inRate)
		var inLen, outLen C.spx_uint32_t
		inLen = C.spx_uint32_t(len(bufIn) / int(C.sizeof_spx_int16_t))
		outLen = C.spx_uint32_t(len(bufOut) / int(C.sizeof_spx_int16_t))
		C.speex_resampler_process_int(rs.instance, C.spx_uint32_t(chanIndex), (*C.spx_int16_t)(unsafe.Pointer(&(bufIn)[0])),
			&inLen, (*C.spx_int16_t)(unsafe.Pointer(&(bufOut)[0])), &outLen)

		bufLen := int(outLen) * int(C.sizeof_spx_int16_t)
		bufOut = bufOut[:bufLen]
	} else {
		bufOut = bufIn
	}

	return
}

func (rs *AudioResampler) Destroy() error {
	if rs.instance != nil {
		C.speex_resampler_destroy(rs.instance)
		rs.instance = nil
		rs.inRate = 0
		rs.outRate = 0
	}
	return nil
}
