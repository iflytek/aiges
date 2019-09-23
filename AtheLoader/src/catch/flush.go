package catch

// TODO 验证极端情况Go层栈数据读取是否异常(C越界写栈异常):内存隔离方案?
// TODO 存储采取S3,同步采取监控系统对接方式.
// TODO 新增线程映射,扩充c层崩溃异常精准化定位程度 (crash threadId && go threadId)

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"time"
)

func catchFlush(rm rpMeta, cStack []byte, goStack []byte) {
	crp := CatchReport{Address: rm.addr, Service: rm.svc, Summary: rm.err.Error(), Pid: rm.pid, Stackc: cStack, Stackgo: goStack}
	crp.Trustlist = nil // TODO trust集合待通过接口上报收集
	if instMgrCallBack != nil {
		doubt := instMgrCallBack(rm.sid)
		// TODO catch analysis
		for _, v := range doubt {
			var md MetaDoubt
			md.Sid = v.Sid
			md.Params = v.Param
			md.Doubt = v.Param // TODO analysis doubtful params
			md.Score = 1
			for _, data := range v.DataList {
				md.Datalist = append(md.Datalist, &MetaData{Data: data.Data, Encoding: data.Enc, Format: data.Fmt, Type: data.Typ})
			}
			crp.Doubtlist = append(crp.Doubtlist, &md)
		}
	}

	cpInfo, err := proto.Marshal(&crp)
	if err != nil {
		catchLog.Errorw("catchFlush pb marshal fail", "error", err.Error())
		return
	}

	if dumpCatch {
		dumpFile := rm.pid + "_" + fmt.Sprintf("%d", time.Now().Unix())
		file, err := os.Create(dumpDir + "/" + dumpFile)
		if err != nil {
			catchLog.Errorw("catchFlush create dump file fail", "error", err.Error())
		} else {
			file.Write(cpInfo)
			file.Close()
		}
	}

	// TODO flush cpInfo to S3 & http notify prometheus
}
