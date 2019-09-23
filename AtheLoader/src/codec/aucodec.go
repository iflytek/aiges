package codec

type codecor interface {
	start(*aucdesc) (interface{}, int, error)
	stop(interface{}) (int, error)
	// 带缓存功能编解码器,encode/decode增加数据状态参数,非整帧尾数据
	encode(interface{}, []byte, bool) ([]byte, int, error)
	decode(interface{}, []byte, bool) ([]byte, int, error)
}

var codecs map[string]codecor

func init() {
	codecs = make(map[string]codecor)
}

type aucdesc struct {
	auc          string // 编解码类型
	compressRate string // 压缩率,不同编解码类型(如speex:0~10)
	sampleRate   string // 采样率,如8000/16000等
	params       string // 编解码参数
	frame        string // 编解码压缩帧大小
}

type AucodecInst struct {
	inst interface{} // 编解码实例
	aucdesc
}

func (ai *AucodecInst) Encode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := codecs[ai.auc]; codec != nil {
		output, code, err = codec.encode(ai.inst, input, last)
		return
	}
	output = input
	return
}

func (ai *AucodecInst) Decode(input []byte, last bool) (output []byte, code int, err error) {
	if codec := codecs[ai.auc]; codec != nil {
		output, code, err = codec.decode(ai.inst, input, last)
		return
	}
	output = input
	return
}

func NewAucodec(tag *map[string]string) (ai *AucodecInst, code int, err error) {
	var desc aucdesc
	if desc, code, err = tagAuc(tag); err == nil {
		ai = &AucodecInst{aucdesc: desc}
		if codec := codecs[ai.auc]; codec != nil {
			ai.inst, code, err = codec.start(&ai.aucdesc)
		}
	}
	return
}

func CloseAucodec(ai *AucodecInst) (code int, err error) {
	if ai != nil {
		if codec := codecs[ai.auc]; codec != nil {
			code, err = codec.stop(ai.inst)
		}
	}
	return
}
