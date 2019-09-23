package protocol

import "github.com/gogo/protobuf/proto"

func rawEncode(in map[string]string) ([]byte, error) {
	ed := EngInputData{}
	ed.EngParam = in
	rst, rstErr := proto.Marshal(&ed)
	return rst, rstErr
}
func RawEncode(in map[string]string) ([]byte, error) {
	return rawEncode(in)
}
