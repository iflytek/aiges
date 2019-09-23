package codec

import (
	"frame"
	"strings"
)

const (
	aucName  = frame.ParaAuCodec
	rateName = frame.ParaAudioRate
	spxName  = frame.ParaSpxFrame
)

// 后续可能迁移至服务实例进行解析
func tagAuc(input *map[string]string) (desc aucdesc, code int, err error) {
	desc.auc = AUDIORAW // default
	if input != nil {
		if auc, exist := (*input)[aucName]; exist {
			desc.auc, desc.compressRate = regularizeAuc(auc)
			desc.sampleRate = (*input)[rateName]
			desc.params = aucParams[desc.auc]
			// adjust speexraw && speexraw-wb
			if fra, exist := (*input)[spxName]; exist {
				desc.frame = fra
				switch desc.auc {
				case AUDIOSPEEX:
					desc.auc = AUDIOSPEEXRAW
				case AUDIOSPEEXWB:
					desc.auc = AUDIOSPEEXRAWWB
				}
			}
			if _, exist := codecs[desc.auc]; !exist {
				// not support audio codec
				code = frame.AigesErrorInvalidParaValue
				err = frame.ErrorAudioCodingStart
			}
		}
	}

	return
}

func regularizeAuc(auc string) (encType string, encRate string) {
	// check audio codec ratio and split it;
	decodec := strings.Split(auc, ";")
	encType = decodec[0]
	encRate = encDefaultRate
	if len(decodec) > 1 {
		encRate = decodec[1]
	}
	return
}
