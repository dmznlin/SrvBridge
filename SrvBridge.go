//添加图标和版本信息
//go:generate goversioninfo -o SrvBridge.syso res/ver.json

package main

import (
  . "github.com/dmznlin/znlib-go/znlib"
)

func main() {
  Info("web server starting")
}
