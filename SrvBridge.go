//添加图标和版本信息
//go:generate goversioninfo -o SrvBridge.syso res/ver.json

package main

import (
	. "SrvBridge/src"
	"github.com/dmznlin/znlib-go/znlib"
)

func main() {
	ServiceWorkGroup.Run(StartThriftService)
	//启动thrift-rpc服务

	znlib.WaitSystemExit(
		func() error { //取消ctx上下文
			ServiceCancel()
			return nil
		},

		StopThriftService, //停止thrift-rpc服务
	) //关闭时清理资源

	ServiceWorkGroup.Wait()
	//等待所有服务退出
}
