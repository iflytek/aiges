package byteutil

import "bytes"

func Merge(items ...[]byte) []byte {
	return bytes.Join(items, []byte(""))
}
