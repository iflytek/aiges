/*
	数据编解码模块,用于处理音频,视频,文本的编解码及压缩/解压操作;
	Note: 开源编解码模块编码解码方式需与客户端约定;
*/
package codec

func AudioCodecInit() (errInfo error) {
	for codec, aucs := range aucMap {
		for _, auc := range aucs {
			switch codec {
			case codecNil:
				codecs[auc] = nil
			case codecSpeex:
				codecs[auc] = &speexCodec{}
			}
		}
	}
	return
}

func AudioCodecFini() {
	return
}
