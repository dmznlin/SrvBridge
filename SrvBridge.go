//添加图标和版本信息
//go:generate goversioninfo -o SrvBridge.syso res/ver.json

package main

import (
	. "SrvBridge/src"
	"fmt"
	"github.com/dmznlin/znlib-go/znlib"
	"github.com/dmznlin/znlib-go/znlib/threading"
)

func main() {
	LoadWorkersConfig()
	//load config first

	var rg *threading.RoutineGroup
	rg = threading.NewRoutineGroup()
	//service works in routine group

	start := func(sw ServiceWorker) {
		rg.RunSafe(func() {
			sw.Start()
		})
	}
	for idx, sw := range ServiceWorkers {
		znlib.Info(fmt.Sprintf("%d.Service [%s] is starting...", idx, sw.WorkName()))
		start(sw)
	} //启动服务

	znlib.WaitSystemExit(
		func() error { //取消ctx上下文
			ServiceCancel()
			return nil
		},

		func() error {
			for idx, sw := range ServiceWorkers {
				znlib.Info(fmt.Sprintf("%d.Service [%s] try to stop...", idx, sw.WorkName()))
				if err := sw.Stop(); err != nil {
					znlib.Error(err.Error())
				}
			}

			return nil
		},
	) //关闭时清理资源

	rg.Wait()
	//等待所有服务退出
}
