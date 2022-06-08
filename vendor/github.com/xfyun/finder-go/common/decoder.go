package common

import (
	"encoding/binary"

	errors "github.com/xfyun/finder-go/errors"
)

func DecodeValue(data []byte) (string, []byte, error) {
	var err error
	if len(data) == 0 {
		err = errors.NewFinderError(errors.ZkGetNilData)
		return "", nil, err
	}
	if len(data) <= 4 {
		err = errors.NewFinderError(errors.InvalidParam)

		return "", nil, err
	}
	l := binary.BigEndian.Uint32(data[:4])
	if int(l) > (len(data) - 4) {
		err = errors.NewFinderError(errors.InvalidParam)
		return "", nil, err
	}
	pushID := string(data[4 : l+4])

	return pushID, data[l+4:], nil
}
