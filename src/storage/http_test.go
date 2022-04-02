package storage

import (
	"fmt"
	"testing"
)

func TestHttpDownload(t *testing.T) {
	data, code, err := HttpDownload("http://google.com")
	if err != nil {
		fmt.Println(code, err)
		return
	}
	fmt.Println(string(data))
}
