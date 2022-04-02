package codec

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var isIgnoreRegexp *regexp.Regexp

type wordInst struct {
	encoding string
}

type wordCodec struct {
}

func (wc *wordCodec) start(it interface{}) (inst interface{}, code int, err error) {
	desc := it.(*txtdesc)
	wdInst := &wordInst{}
	wdInst.encoding = desc.encoding
	isIgnoreRegexp, _ = regexp.Compile(`^[A-Za-z!#$%&()*+=_]+$`)
	return wdInst, 0, nil
}

func (wc *wordCodec) stop(inst interface{}) (code int, err error) {
	return
}

func (wc *wordCodec) encode(inst interface{}, input []byte, last bool) (output []byte, code int, err error) {
	return input, 0, nil
}

//当前是做校验 用的，看看传输的数据和设置的文本编码是否相符
func (wc *wordCodec) decode(inst interface{}, input []byte, last bool) (output []byte, code int, err error) {
	wdInst := inst.(*wordInst)
	switch wdInst.encoding {

	case TEXTUTF8:
		if utf8.Valid(input) {
			return input, 0, nil
		} else {
			return nil, -1, errors.New("content is inconsistent with the encoding")
		}
	case TEXTGB2312:
		if isIgnoreRegexp.MatchString(string(input)) {
			return input, 0, nil
		}
		if utf8.Valid(input) {
			return nil, -1, errors.New("content is inconsistent with the encoding")
		}
		return input, 0, nil
	default:
		return input, 0, nil
	}
}
