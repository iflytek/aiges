/*
* @file	roundbinlb.go
* @brief  native 无权重负载均衡的实现
*
* @author	kunliu2
* @version	1.0
* @date		2017.12
*/
package xsf

type hashLB struct {
	sd *serviceDiscovery
}

func newHashLB(o *conOption) *hashLB {
	hashl := new(hashLB)
	hashl.sd = newServiceDiscovery(o.fm)
	return hashl
}
func (l *hashLB) strategyJudge(param *LBParams) bool {
	if "" != param.hashKey {
		return true
	}
	return false
}

func (l *hashLB) getNbest(addrs []string, ix, nbest int) []string {
	addrsLen := len(addrs)

	if 0 == addrsLen && nbest <= 0 && ix < 0 {
		return nil
	}

	if addrsLen < nbest {
		return addrs
	}

	var rst []string
	if ix+nbest > addrsLen {
		rst = addrs[ix:]
		if t := nbest - len(rst); t > 0 {
			rst = append(rst, addrs[:t]...)
		}
	} else {
		rst = addrs[ix : ix+nbest]
	}

	return rst
}

func (l *hashLB) Find(param *LBParams) ([]string,[]string, error) {

	/*
		1、如果有hashkey，采用hash策略，否则轮询

	*/
	var s *service
	var e error
	if l.strategyJudge(param) {
		//采用hash策略
		param.log.Infow("get hashKey",
			"fn", "Find", "strategy", "hash", "hashKey", param.hashKey)

		s, e = l.sd.findAll(param.version, param.svc, param.logId, param.log)
		if e == nil {

			param.log.Infow("l.sd.findAll success",
				"logId", param.logId, "s", s, "e", e)

			//此处效率待优化
			addrsTmp,_ := s.addrs.NextInList(0)
			addrsTmpLen := len(addrsTmp)
			ix := l.handle2Ip(param.hashKey, addrsTmpLen)

			addrs := l.getNbest(addrsTmp, ix, param.nbest)
			param.log.Infow("hashLb getNbest",
				"fn", "getNbest", "addrs", addrs, "addrsTmp", addrsTmp, "ix", ix, "nbest", param.nbest, "logId", param.logId)


			if len(addrs) > 0 {
				return addrs,nil, nil
			} else {
				param.log.Infow("can't take enough addrs", "logId", param.logId, "s", s, "e", e)
				return nil,nil, EINVALIDLADDR
			}

		}

	} else {

		//	退化为轮询策略
		s, e = l.sd.findAll(param.version, param.svc, param.logId, param.log)
		if e == nil {

			param.log.Infow("l.sd.findAll success", "logId", param.logId, "s", s, "e", e)

			addrs,allAddrs := s.addrs.NextInList(param.nbest)

			param.log.Infow("fn:Find", "adds", addrs, "logId", param.logId)

			if len(addrs) > 0 {
				return addrs,allAddrs, nil
			} else {
				param.log.Infow("can't take enough addrs", "logId", param.logId, "s", s, "e", e)
				return nil,nil, EINVALIDLADDR
			}
		}
	}

	param.log.Infow("l.sd.findAll failed", "logId", param.logId, "s", s, "e", e)

	return nil,nil, e
}

func (l *hashLB) handle2Ip(handle string, seed int) (ix int) {
	var cnt int
	for ix := 0; ix < len(handle); ix++ {
		cnt += int(handle[ix])
	}
	return cnt % seed
}
