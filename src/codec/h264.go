package codec

/*
#cgo CFLAGS: -I../cgo/header/ffmpeg
#cgo LDFLAGS: -L../cgo/library -lavcodec -lavutil -lswresample
#include <stdlib.h>
#include "./libavcodec/avcodec.h"
#include "./libavutil/imgutils.h"
*/
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

func ffmpegInit() {
	// C.avcodec_register_all()
}

type h264Inst struct {
	av  *C.AVCodec        // 解码器
	cxt *C.AVCodecContext // 解码上下文
	frm *C.AVFrame        // 解码后二维图像
	yuv *C.AVFrame
	dic *C.AVDictionary
	pck C.AVPacket
}

type h264Codec struct{}

func (hc *h264Codec) start(it interface{}) (inst interface{}, code int, err error) {
	desc := it.(*vdcdesc)
	// TODO 异常场景逆初始化
	h264 := &h264Inst{}
	h264.av = C.avcodec_find_decoder(C.AV_CODEC_ID_H264)
	if h264.av == nil {
		return nil, -1, errors.New("avcodec_find_decoder decoder can't find")
	}

	h264.cxt = C.avcodec_alloc_context3(h264.av)
	if h264.cxt == nil {
		return nil, -1, errors.New("avcodec_alloc_context3 alloc context fail")
	}
	// context 参数设置
	(*h264.cxt).codec_type = C.AVMEDIA_TYPE_VIDEO
	(*h264.cxt).height = C.int(desc.heigth)
	(*h264.cxt).width = C.int(desc.width)

	ret := C.avcodec_open2(h264.cxt, h264.av, (**C.AVDictionary)(unsafe.Pointer(&h264.dic)))
	if ret < 0 {
		C.avcodec_close(h264.cxt)
		C.avcodec_free_context(&h264.cxt)
		return nil, int(ret), errors.New("avcodec_open2 open fail")
	}

	C.av_init_packet(&h264.pck)
	h264.frm = C.av_frame_alloc()
	h264.yuv = C.av_frame_alloc()

	inst = h264
	return
}

func (hc *h264Codec) stop(inst interface{}) (code int, err error) {
	h264 := inst.(*h264Inst)
	C.av_packet_unref(&h264.pck)
	C.av_frame_free(&h264.frm)
	C.av_frame_free(&h264.yuv)
	C.avcodec_close(h264.cxt)
	C.avcodec_free_context(&h264.cxt)
	return
}

func (hc *h264Codec) encode(inst interface{}, input []byte, last bool) (output []byte, code int, err error) {
	// TODO h264编码
	return
}

func (hc *h264Codec) decode(inst interface{}, input []byte, last bool) (output []byte, code int, err error) {
	if len(input) == 0 {
		return nil, -1, errors.New("the video is empty")
	}
	h264 := inst.(*h264Inst)
	h264.pck.data = (*C.uchar)(unsafe.Pointer(&input[0]))
	h264.pck.size = C.int(len(input))
	var picFrame C.int
	// TODO 非整帧缓存,整帧解码.(新增帧头读取处理, 多帧写入导致异常)
	ret := C.avcodec_decode_video2(h264.cxt, h264.frm, &picFrame, &h264.pck)
	if ret < 0 {
		return nil, int(ret), errors.New("avcodec_decode_video2 fail with " + strconv.Itoa(int(ret)))
	} else if ret != h264.pck.size {
		// 整帧解码,写入数据全部解码,上层缓存非整帧数据
		return nil, -1, errors.New("number of bytes used incorrect")
	} else {
		// 获取解码图像
		if picFrame == 0 {
			return // 相关视频编码帧不带实体数据
			// return nil, -1, errors.New("no frame could be decompressed")
		}
		output, code, err = hc.frame2pic(h264)
		C.av_packet_unref(&h264.pck)
		C.av_init_packet(&h264.pck)
	}
	return
}

// 获取解码图像数据及图像处理
func (hc *h264Codec) frame2pic(h264 *h264Inst) (output []byte, code int, err error) {
	// TODO 图像转换; 缺省AV_PIX_FMT_YUV420P,(见:enum AVPixelFormat)

	buffSize := C.av_image_get_buffer_size((*h264.cxt).pix_fmt, (*h264.cxt).width, (*h264.cxt).height, 16)
	buff := C.av_malloc(C.ulong(buffSize))
	defer C.av_freep(unsafe.Pointer(&buff))
	writes := C.av_image_copy_to_buffer((*C.uchar)(buff), buffSize, (**C.uint8_t)(&h264.frm.data[0]), (*C.int)(&h264.frm.linesize[0]),
		h264.cxt.pix_fmt, h264.cxt.width, h264.cxt.height, 1)
	if writes < 0 {
		return nil, int(writes), errors.New("av_image_copy_to_buffer fail")
	}
	output = C.GoBytes(unsafe.Pointer(buff), writes)
	return
}
