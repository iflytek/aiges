/*
 *实现GetServer接口
 */
package daemon

import (
	"git.xfyun.cn/AIaaS/xsf-external/server"
	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

type StrategyInst interface {
	setServer(...SetInPutOpt) (err LbErr)
	getServer(...GetInPutOpt) (nBestNodes []string, nBestNodesErr LbErr)
	init(toolbox *xsf.ToolBox)
	serve(in *xsf.Req, span *xsf.Span, toolbox *xsf.ToolBox) (res *utils.Res, err error)
}
type LbHandle struct {
	toolbox  *xsf.ToolBox
	worker   StrategyInst
	strategy StrategyClassify
}

func (lh *LbHandle) Init(toolbox *xsf.ToolBox) (err error) {
	lh.toolbox = toolbox
	strategy, err := lh.toolbox.Cfg.GetInt(BO, STRATEGY)
	if nil != err {
		lh.toolbox.Log.Errorf("toolbox parse config param strategy error:%s", err.Error())
		return
	}

	std.Printf("strategy:%v\n", strategy)
	lh.toolbox.Log.Errorf("strategy:%v", strategy)
	lh.strategy = StrategyClassify(strategy)
	if lh.strategy.String() == Unknown {
		return ErrLbStrategyIsNotSupport
	}
	switch lh.strategy {
	case load:
		{
			std.Println("about to call newLoad")
			lh.worker = newLoad(lh.toolbox)
		}
	case poll:
		{
			std.Printf("about to call newPoll\n")
			lh.worker = newPoll(lh.toolbox)
		}
	case loadMini:
		{
			std.Printf("about to call newLoadMini\n")
			lh.worker = newLoadMini(lh.toolbox)
		}
	}

	lh.toolbox.Log.Errorf("start strategy:%v error:%v", strategy, err)

	return

}

func (lh *LbHandle) Stop() (err error) {
	return
}
