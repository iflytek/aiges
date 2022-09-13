package codec

const (
	decMaxPerFrame    int    = 640 // 极限场景,编解码压缩比达到640:1
	encRatePrefix     string = "mode="
	encDefaultRate    string = "7"
	iflyIcoParams     string = "bandwidth=7000,bitrate=16000"
	iflyLameParam     string = "samplerate=16000"
	iflyLame8kParam   string = "samplerate=8000"
	iflyOpusH8Param   string = "samplerate=8000"
	iflyOpusWbH8Param string = "samplerate=16000"
	iflyOpusOggParam  string = "samplerate=16000"
	iflyCodecLib      string = "libspeex.so;libico.so;libopus.so;libict.so;liblame.so;libamr.so;libamr_wb.so"
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
	AUDIOICT        string = "ict"
	AUDIOOPUS       string = "opus"
	AUDIOOPUSWB     string = "opus-wb"
	AUDIOOPUSSWB    string = "opus-swb"
	AUDIOOPUSOGG    string = "opus-ogg"
	AUDIOOPUSH8     string = "opus-h8"
	AUDIOOPUSWBH8   string = "opus-wb-h8"
	AUDIOLAME       string = "lame"
	AUDIOSPEEXRAW   string = "speex-org-nb" // 开源speex narrowband
	AUDIOSPEEXRAWWB string = "speex-org-wb" // 开源speex wideband
)

// 加载器支持的音频编解码器种类
const (
	codecNil = iota
	codecSpeex
)

// 各编解码器对应支持的编解码类型
var aucMap = map[int][]string{
	codecNil:   []string{AUDIORAW},
	codecSpeex: []string{AUDIOSPEEXRAW, AUDIOSPEEXRAWWB},
}

// 加载器支持的文本编解码器
const (
	codecTxtNil = iota
	codecWord
)

const (
	TEXTRAW    string = ""
	TEXTUTF8   string = "utf8"
	TEXTGB2312 string = "gb2312"
)

var txtMap = map[int][]string{
	codecTxtNil: []string{TEXTRAW},
	codecWord:   []string{TEXTUTF8, TEXTGB2312},
}

// 加载器支持的图像编解码器
const (
	codecImgNil = iota
)

const (
	IMGRAW string = ""
)

var imgMap = map[int][]string{
	codecImgNil: []string{IMGRAW},
}

// 加载器支持的视频编解码器
const (
	codecVdNil = iota
	codecH264
)

const (
	VIDEORAW  string = ""
	VIDEOH264 string = "h264"
)

var vdcMap = map[int][]string{
	codecVdNil: []string{VIDEORAW},
	codecH264:  []string{},
}
