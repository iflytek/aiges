package codec

const (
	decMaxPerFrame int    = 640 // 极限场景,编解码压缩比达到640:1
	encRatePrefix  string = "mode="
	encDefaultRate string = "7"
	iflyIcoParams  string = "bandwidth=7000,bitrate=16000"
	iflyCodecLib   string = "libamr.so;libamr_wb.so;libspeex.so;libico.so;libopus.so;libict.so;liblame.so"
)

// 云端所支持编解码类型;
const (
	AUDIORAW        string = "raw"
	AUDIOSPEEX      string = "speex"
	AUDIOSPEEXWB    string = "speex-wb"
	AUDIOAMR        string = "amr"
	AUDIOAMRWB      string = "amr-wb"
	AUDIOAMRWBFX    string = "amr-wb-fx"
	AUDIOICO        string = "ico"
	AUDIOOPUS       string = "opus"
	AUDIOOPUSWB     string = "opus-wb"
	AUDIOOPUSOGG    string = "opus-ogg"
	AUDIOLAME		string = "lame"
	AUDIOSPEEXRAW   string = "speexraw"
	AUDIOSPEEXRAWWB string = "speexraw-wb"
)

// 加载器支持的编解码器种类
const (
	codecNil = iota
	codecSpeex
)

// 各编解码器对应支持的编解码类型
var aucMap = map[int][]string{
	codecNil:   []string{AUDIORAW},
	codecSpeex: []string{AUDIOSPEEXRAW, AUDIOSPEEXRAWWB},
}

var aucParams = map[string]string{
	AUDIOICO: iflyIcoParams,
}
