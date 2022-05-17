//添加图标和版本信息
//go:generate goversioninfo -o SrvBridge.syso res/ver.json

package main

import (
	. "SrvBridge/src"
)

func main() {
	StartThriftService()
	//启动thrift-rpc服务
}
