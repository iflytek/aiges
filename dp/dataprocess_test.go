//go:build linux
// +build linux

package dp

import (
	"bytes"
	"io/ioutil"
	"testing"
)

/*
	1. 8k升采样
	2. 16k降采样
*/

const (
	audioFile8K   = "../../test/audio/test8k.pcm"
	audioFile16K  = "../../test/audio/test.pcm"
	audioUpFile   = "../../test/audio/test8kTo16k.pcm"
	audioDownFile = "../../test/audio/testTo8k.pcm"
)

func Test8kTo16k(t *testing.T) {
	ar := AudioResampler{}
	// read audio slice
	audio, err := ioutil.ReadFile(audioFile8K)
	if err != nil {
		t.Errorf("read audio file %s fail with %s", audioFile8K, err.Error())
	}
	audio16k, err := ioutil.ReadFile(audioUpFile)
	if err != nil {
		t.Errorf("read audio file %s fail with %s", audioUpFile, err.Error())
	}

	_ = ar.Init(1, 8000, 16000, 10)
	audioBuf, err := ar.ProcessInt(0, audio)
	if err != nil {
		t.Errorf("resample audio fail with %s", err.Error())
	}

	// 与本地音频校验check
	ret := bytes.Compare(audioBuf, audio16k)
	if ret != 0 {
		t.Errorf("resample src audio %s is not invalid, dst audio %s", audioFile8K, audioUpFile)
	}

	_ = ar.Destroy()
}

func Test16kTo8k(t *testing.T) {
	ar := AudioResampler{}
	// read audio slice
	audio, err := ioutil.ReadFile(audioFile16K)
	if err != nil {
		t.Errorf("read audio file %s fail with %s", audioFile16K, err.Error())
	}
	audio8k, err := ioutil.ReadFile(audioDownFile)
	if err != nil {
		t.Errorf("read audio file %s fail with %s", audioDownFile, err.Error())
	}

	_ = ar.Init(1, 16000, 8000, 10)
	audioBuf, err := ar.ProcessInt(0, audio)
	if err != nil {
		t.Errorf("resample audio fail with %s", err.Error())
	}

	// 与本地音频校验check
	ret := bytes.Compare(audioBuf, audio8k)
	if ret != 0 {
		t.Errorf("resample src audio %s is not invalid, dst audio %s", audioFile16K, audioDownFile)
	}

	_ = ar.Destroy()
}
