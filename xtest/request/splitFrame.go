package request

/*
#cgo CFLAGS: -I../include
#cgo LDFLAGS: -L../lib -lh264bitstream
#include "stdlib.h"
#include <stdint.h>
#include "stdio.h"
#include "h264_nalu_split.h"
*/
import "C"
import "unsafe"


func GetH264Frames(video []byte) (frameSizes []int) {
	videoLen := len(video)
	naluList := C.get_h264_nalu((*C.uchar)(unsafe.Pointer(&video[0])), C.ulong(videoLen))
	tmp:=0
	for naluList != nil {
		//仅仅获取图片帧
		if int(naluList.nalu_type)==7 || int(naluList.nalu_type)==8 {
			tmp+=int(naluList.start_code_len+naluList.size)
		}else{
			frameSizes=append(frameSizes,int(naluList.start_code_len+naluList.size)+tmp)
			if tmp!=0 {
				tmp=0
			}
		}
		naluList=naluList.next
	}
	return
}
