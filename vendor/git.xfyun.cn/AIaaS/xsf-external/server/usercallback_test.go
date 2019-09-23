package xsf

import (
	"fmt"
	"testing"
)

type userCallBack struct {
	d interface{}
}

func (e *userCallBack) Exec() {
	fmt.Println(e.d)
}
func Test_dealUserCallBack(t *testing.T) {
	{
		fmt.Println(AddUserCallBack("xxx", &userCallBack{"xxx"}))
		fmt.Println(AddUserCallBack("xxx", &userCallBack{"xxx"}))
	}

	{
		fmt.Println(AddUserCallBackWithPriority("UserLowPriority", &userCallBack{"UserLowPriority"}, UserLowPriority), 1)
		fmt.Println(AddUserCallBackWithPriority("UserNormalPriority", &userCallBack{"UserNormalPriority"}, UserNormalPriority), 2)
		fmt.Println(AddUserCallBackWithPriority("UserHighPriority", &userCallBack{"UserHighPriority"}, UserHighPriority), 3)
	}

	{
		dealUserCallBack()
	}
}
func Test_deal(t *testing.T) {

	{
		dealUserCallBack()
	}
}
