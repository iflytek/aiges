package ratelimit

import "sync/atomic"

var(
	tokenNum  int32 = 0;
)

func GetToken()int32  {

	return tokenNum
}

func ReleaseToken()  {
	atomic.StoreInt32(&tokenNum,tokenNum-1)
}