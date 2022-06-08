package codec

import "fmt"

// 编解码器
// encode/decode增加数据状态参数,处理非整帧尾数据
type codecor interface {
	start(interface{}) (interface{}, int, error)
	stop(interface{}) (int, error)
	encode(interface{}, []byte, bool) ([]byte, int, error)
	decode(interface{}, []byte, bool) ([]byte, int, error)
}

// 编解码实例
type Instance interface {
	Encode(input []byte, last bool) (output []byte, code int, err error)
	Decode(input []byte, last bool) (output []byte, code int, err error)
}

func CodecInit() (err error) {
	if err = auCodecInit(); err != nil {
		return
	}
	if err = vdCodecInit(); err != nil {
		return
	}
	if err = txtCodecInit(); err != nil {
		return
	}
	if err = imgCodecInit(); err != nil {
		return
	}
	return
}

func CodecFini() {
	auCodecFini()
	fmt.Println("aiService.Finit: fini codeCFini success!")
}

// 音频编解码初始化
func auCodecInit() (err error) {
	// raw speex codecor no need init
	for codec, aucs := range aucMap {
		for _, auc := range aucs {
			switch codec {
			case codecNil:
				aucodecs[auc] = nil
			case codecSpeex:
				aucodecs[auc] = &speexCodec{}
			}
		}
	}
	return
}

func auCodecFini() {
	return
}

// 视频编解码初始化
func vdCodecInit() (err error) {
	for codec, vdcs := range vdcMap {
		for _, vdc := range vdcs {
			switch codec {
			case codecVdNil:
				vdcodecs[vdc] = nil
			case codecH264:
				vdcodecs[vdc] = &h264Codec{}
			}
		}
	}
	return
}

// 文本编解码初始化
func txtCodecInit() (err error) {
	for codec, txtcs := range txtMap {
		for _, txtc := range txtcs {
			switch codec {
			case codecTxtNil:
				txtcodecs[txtc] = nil
			case codecWord:
				txtcodecs[txtc] = &wordCodec{}
			}
		}
	}
	return
}

// 图像编解码初始化
func imgCodecInit() (err error) {
	for codec, imgcs := range imgMap {
		for _, imgc := range imgcs {
			switch codec {
			case codecImgNil:
				imgcodecs[imgc] = nil
			}
		}
	}
	return
}
