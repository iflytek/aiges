/*
* @file	lbi.go
* @brief  负载均衡interface
*
* @author	kunliu2
* @version	1.0
* @date		2017.12
 */

package xsf

import (
	"fmt"
)

// LoadBalce 的模式
type LBMode int

// String LBMode的串化
func (s LBMode) String() string {
	switch s {
	case RoundRobin:
		return "RoundRobin"
	case Hash:
		return "Hash"
	case Remote:
		return "Remote"
	case ConHash:
		return "ConHash"
	case TopKHash:
		return "TopKHash"
	default:
		return "Invalid-LBMode"
	}
}

// LBMode取值范围
const (
	//  轮循
	RoundRobin LBMode = iota
	//  hash
	Hash
	// 远端负载
	Remote
	//  一致性hash
	ConHash
	// topK
	TopKHash
)

// LBParams 传给LB的相关参数
type LBParams struct {
	version string

	// name LB的名字
	name string

	// svc 服务名称
	svc string

	// catgory业务名
	catgory string

	// ext 扩展参数
	ext map[string]string

	// nbest 获取的机器个数
	nbest int

	//retry lb 重试次数
	try int

	//日志相关
	log  *Logger
	span *Span

	logId string

	hashKey string

	retry bool //表示本次为重试请求

	localIp string

	directEngIp string

	peerIp string

	failed string //上次失败的节点
}

func (p *LBParams) String() string {
	return fmt.Sprintf("name:%v,svc:%v,catgory:%v,nbest:%v,try:%v,directEngIp:%v", p.name, p.svc, p.catgory, p.nbest, p.try, p.directEngIp)
}

func (p *LBParams) WithLocalIp(localIp string) {
	p.localIp = localIp
}

func (p *LBParams) WithVersion(version string) {
	p.version = version
}

// WithName 设置LB服务名
func (p *LBParams) WithName(name string) {
	p.name = name
}

// WithSvc 设置服务的名字
func (p *LBParams) WithSvc(svc string) {
	p.svc = svc
}

func (p *LBParams) WithLogId(logId string) {
	p.logId = logId
}

func (p *LBParams) WithDirectEngIp(directEngIp string) {
	p.directEngIp = directEngIp
}
func (p *LBParams) WithPeerIp(PeerIp string) {
	p.peerIp = PeerIp
}

func (p *LBParams) WithRetry(retry bool) {
	p.retry = retry
}

// WithCatgory 设置业务类别
func (p *LBParams) WithCatgory(c string) {
	p.catgory = c
}

// WithExtend 设置扩展的参数
func (p *LBParams) WithExtend(ext map[string]string) {
	p.ext = ext
}

// WithLog 设置日志操作句柄
func (p *LBParams) WithLog(log *Logger) {
	p.log = log
}

// WithTracer 设置trace句柄
func (p *LBParams) WithTracer(sp *Span) {
	p.span = sp
}

// WithNBest 设置trace句柄
func (p *LBParams) WithNBest(n int) {
	if n == 0 && p.retry {
		p.nbest = 1
		return
	}
	p.nbest = n
}

func (p *LBParams) WithHashKey(key string) {
	p.hashKey = key
}
func (p *LBParams) WithFailed(failed string) {
	p.failed = failed
}

/*
// WithTracer 设置trace句柄
func (p *LBParams)WithTry( n int32){
	if n == 0{
		p.try = 2
		return
	}
	p.try = n
}
*/

// 负责均衡实际实现的interface
type LBI interface {
	Find(param *LBParams) ([]string, []string, error) //nbest addr,all addr
}
