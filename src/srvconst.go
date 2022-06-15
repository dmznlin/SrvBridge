/*Package srvbridge ***********************************************************
  作者: dmzn@163.com 2022-06-09 11:55:04
  描述: 公用常量变量定义
******************************************************************************/
package srvbridge

import (
	"context"
	"github.com/dmznlin/znlib-go/znlib"
	"sync"
)

type ServiceWorker interface {
	WorkName() string //服务名称
	Start()           //启动函数
	Stop() error      //停止函数
}

var (
	global_init    sync.Once       //全局初始化
	ServiceWorkers []ServiceWorker //服务对象列表

	ServiceContext context.Context    //服务上下文
	ServiceCancel  context.CancelFunc //取消上下文
)

func init() {
	Init_global()
}

func Init_global() {
	global_init.Do(func() {
		ServiceWorkers = make([]ServiceWorker, 0, 5)
		ServiceContext, ServiceCancel = context.WithCancel(context.Background())
	})
}

/*RegistWorker 2022-06-10 10:33:08
  参数: sw,服务对象
  描述: 注册一个服务对象
*/
func RegistWorker(sw ServiceWorker) {
	znlib.Application.SyncLock.Lock()
	defer znlib.Application.SyncLock.Unlock()

	Init_global()
	ServiceWorkers = append(ServiceWorkers, sw)
}
