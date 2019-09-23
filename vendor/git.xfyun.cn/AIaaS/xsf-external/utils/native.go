package utils

import "io/ioutil"

/*
nativeCfg:
读取本地文件
*/
type nativeCfg struct{}

func (*nativeCfg) Read(name string) (string, error) {
	c, e := ioutil.ReadFile(name)
	return string(c), e
}

/*NewNative:
创建本地配置读取interface
*/
func NewNative(co *CfgOption) (*Configure, error) {
	c := new(Configure)
	e := c.init(new(nativeCfg), co)
	if nil != e {
		return nil, e
	}
	return c, e
}

type bytesReader struct {
	cfgData string
}

func (b *bytesReader) Read(_ string) (string, error) {
	return b.cfgData, nil
}
func newBytesReader(cfgData string) (*Configure, error) {
	c := new(Configure)
	e := c.init(
		&bytesReader{cfgData: cfgData},
		&CfgOption{})
	if nil != e {
		return nil, e
	}
	return c, e
}
