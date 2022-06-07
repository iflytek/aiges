package protocol

import (
	"github.com/golang/protobuf/proto"
	"github.com/xfyun/aiges/frame"
)

const (
	ReqSrc = "source"
	SrcV2  = "aipaas"
)

func InputAdapter(input []byte, ei *LoaderInput) (code int, err error) {
	err = proto.Unmarshal(input, ei)
	if err != nil {
		return frame.AigesErrorPbUnmarshal, err
	}
	if ei.Headers == nil {
		ei.Headers = make(map[string]string)
	}
	ei.Headers[ReqSrc] = SrcV2
	return
}

func OutputAdapter(eo *LoaderOutput) (output []byte, code int, err error) {
	output, err = proto.Marshal(eo)
	if err != nil {
		return nil, frame.AigesErrorPbMarshal, frame.ErrorPbMarshal
	}
	return
}
