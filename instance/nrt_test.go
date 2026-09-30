package instance

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/xfyun/aiges/buffer"
	"github.com/xfyun/aiges/frame"
	"github.com/xfyun/aiges/protocol"
)

func TestNrtDataFillRejectsHTTPWithoutFetching(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		_, _ = w.Write([]byte("private sentinel"))
	}))
	defer server.Close()

	for _, url := range []string{server.URL, server.URL + "/redirect", "https://127.0.0.1/private", "not-a-url", ""} {
		t.Run(url, func(t *testing.T) {
			original := []byte("inline input")
			data := []buffer.DataMeta{{
				Data: original,
				Desc: &protocol.MetaDesc{Attribute: map[string]string{"dataSrc": "http", "url": url}},
			}}
			code, err := nrtDataFill(&data)
			if code != frame.AigesErrorInvalidData || err == nil {
				t.Fatalf("expected invalid input error, got code=%d err=%v", code, err)
			}
			if err.Error() != "http input data source is no longer supported; send data directly" {
				t.Fatalf("unexpected public error: %v", err)
			}
			if !bytes.Equal(data[0].Data.([]byte), original) {
				t.Fatal("rejected request changed input data")
			}
		})
	}
	if got := atomic.LoadInt32(&requests); got != 0 {
		t.Fatalf("rejected HTTP input issued %d requests", got)
	}
}

func TestNrtDataFillPreservesInlineData(t *testing.T) {
	for _, attrs := range []map[string]string{nil, {"dataSrc": "client"}} {
		original := []byte("direct input")
		data := []buffer.DataMeta{{Data: original, Desc: &protocol.MetaDesc{Attribute: attrs}}}
		if code, err := nrtDataFill(&data); code != 0 || err != nil {
			t.Fatalf("inline input rejected: code=%d err=%v", code, err)
		}
		if !bytes.Equal(data[0].Data.([]byte), original) {
			t.Fatal("inline input changed")
		}
	}
}

func TestNrtDataFillRejectsHTTPInMixedBatch(t *testing.T) {
	data := []buffer.DataMeta{
		{Data: []byte("direct input"), Desc: &protocol.MetaDesc{Attribute: map[string]string{"dataSrc": "client"}}},
		{Desc: &protocol.MetaDesc{Attribute: map[string]string{"dataSrc": "http", "url": "http://127.0.0.1/private"}}},
	}
	if code, err := nrtDataFill(&data); code != frame.AigesErrorInvalidData || err == nil {
		t.Fatalf("mixed batch accepted: code=%d err=%v", code, err)
	}
}
