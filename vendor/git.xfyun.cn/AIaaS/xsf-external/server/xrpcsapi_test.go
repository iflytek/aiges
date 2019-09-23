package xsf

import (
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"sync"
	"testing"
)

type userserver struct {
	tool *ToolBox
}

func (c *userserver) Init(toolbox *ToolBox) error {
	c.tool = toolbox
	c.tool.Log.Errorw("this is a error log msg from userserver1.", "func", "Init")
	return nil
}
func (c *userserver) FInit() error {
	c.tool.Log.Errorw("this is a error log msg from userserver1.", "func", "FInit")
	return nil
}
func (c *userserver) Call(in *Req, span *Span) (*Res, error) {
	c.tool.Log.Errorw("this is a error log msg from userserver1.", "func", "Call")
	res := NewRes()
	return res, nil
}

func TestRun(t *testing.T) {
	var server1 XsfServer
	bc := BootConfig{CfgMode: utils.Native,
		CfgData: CfgMeta{CfgName: "test.toml", Project: "test", Group: "default", Service: "xsf", Version: "1.0.0", CompanionUrl: "http://10.1.86.228:9080"}}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server1.Run(bc, &userserver{}); err != nil {
			t.Fatal(err)
		}
	}()
	wg.Wait()
}
