package codec

import (
	"github.com/xfyun/aiges/frame"
	"strconv"
	"strings"
)

type TextTag struct {
	Encoding string
	Compress string
}

type ImageTag struct {
	Encoding string
}

type AudioTag struct {
	Encoding   string
	SampleRate string
	SpeexFrame string
}

type VideoTag struct {
	Encoding  string
	FrameRate string
	Width     string
	Height    string
}

// 后续可能迁移至服务实例进行解析
func normalize(at AudioTag) (desc aucdesc, code int, err error) {
	desc.auc = AUDIORAW // default
	if len(at.Encoding) > 0 {
		desc.auc, desc.compressRate = regularizeAuc(at.Encoding)
		desc.sampleRate = at.SampleRate
		// adjust speexraw && speexraw-wb
		spxFsize, _ := strconv.Atoi(at.SpeexFrame)
		if spxFsize > 0 {
			desc.frame = at.SpeexFrame
			switch desc.auc {
			case AUDIOSPEEX:
				desc.auc = AUDIOSPEEXRAW
			case AUDIOSPEEXWB:
				desc.auc = AUDIOSPEEXRAWWB
			}
		}
		if _, exist := aucodecs[desc.auc]; !exist {
			// not support audio codec
			code = frame.AigesErrorInvalidParaValue
			err = frame.ErrorAudioCodingStart
		}
	}

	return
}

func regularizeAuc(aue string) (encType string, encRate string) {
	// check audio codec ratio and split it;
	decodec := strings.Split(aue, ";")
	encType = decodec[0]
	encRate = encDefaultRate
	if len(decodec) > 1 {
		encRate = decodec[1]
	}
	return
}
