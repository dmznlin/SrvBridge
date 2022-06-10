/*Package srvbridge ***********************************************************
  作者: dmzn@163.com 2022-06-09 11:55:04
  描述: 公用常量变量定义
******************************************************************************/
package srvbridge

import (
	"context"
	"sync"
)

type ServiceWorker interface {
	WorkName() string //服务名称
	Start()           //启动函数
	Stop() error      //停止函数
}

var (
	SyncLock       sync.RWMutex    //全局同步锁
	ServiceWorkers []ServiceWorker //服务对象列表

	ServiceContext context.Context    //服务上下文
	ServiceCancel  context.CancelFunc //取消上下文
)

func init() {
	ServiceWorkers = make([]ServiceWorker, 0, 5)
	ServiceContext, ServiceCancel = context.WithCancel(context.Background())
}

/*RegisteWorker 2022-06-10 10:33:08
  参数: sw,服务对象
  描述: 注册一个服务对象
*/
func RegisteWorker(sw ServiceWorker) {
	SyncLock.Lock()
	defer SyncLock.Unlock()
	ServiceWorkers = append(ServiceWorkers, sw)
}
