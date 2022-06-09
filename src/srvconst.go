/*Package srvbridge ***********************************************************
  作者: dmzn@163.com 2022-06-09 11:55:04
  描述: 公用常量变量定义
******************************************************************************/
package srvbridge

import (
	"context"
	"github.com/dmznlin/znlib-go/znlib/threading"
)

var (
	ServiceWorkGroup *threading.RoutineGroup //服务组
	ServiceContext   context.Context         //服务上下文
	ServiceCancel    context.CancelFunc      //取消上下文
)

func init() {
	ServiceWorkGroup = threading.NewRoutineGroup()
	ServiceContext, ServiceCancel = context.WithCancel(context.Background())
}
