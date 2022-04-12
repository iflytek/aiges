/*
* @file	roundbinlb.go
* @brief  native 无权重负载均衡的实现
*
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */
package xsf

type rrLB struct {
	sd *serviceDiscovery
}

func newRRLB(o *conOption) *rrLB {
	rrl := new(rrLB)
	rrl.sd = newServiceDiscovery(o.fm)
	return rrl
}

/*
func (l *rrLB)Find(param *LBParams)([]string, error){
	s,e:=l.sd.findAll(param.svc)

	if e == nil {
		r:=s.addrs.Next()
		if  r ==  nil {
			return nil, EINVALIDADDR
		}
		addrs:=make([]string,0, 1)
		addrs = append(addrs,r.addr)
		return  addrs,nil
	}
	return nil, e
}*/

func (l *rrLB) Find(param *LBParams) ([]string, []string, error) {
	s, e := l.sd.findAll(param.version, param.svc, param.logId, param.log)
	if e == nil {

		param.log.Infow("l.sd.findAll success", "logId", param.logId, "s", s, "e", e)

		addrs, allAddrs := s.addrs.NextInList(param.nbest)
		param.log.Infow("fn:Find", "adds", addrs)

		if len(addrs) > 0 {
			return addrs, allAddrs, nil
		} else {
			param.log.Infow("can't take enough addrs", "logId", param.logId, "s", s, "e", e)
			return nil, nil, EINVALIDLADDR
		}
	}

	param.log.Infow("l.sd.findAll failed", "logId", param.logId, "s", s, "e", e)

	return nil, nil, e
}

/*
func newRRLB(l LBI,o ...connOpt) *loadBalance{
	lb := new(loadBalance)
	lb.lbi = l
	lb.conns = newConnPool(o...)// 连接属性
	return lb
}
*/
