/*Package srvbridge ***********************************************************
  作者: dmzn@163.com 2022-06-09 13:03:31
  描述: service for thrift
******************************************************************************/
package srvbridge

import (
	. "SrvBridge/src/mit"
	"context"
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
	. "github.com/dmznlin/znlib-go/znlib"
	inifile "github.com/go-ini/ini"
	"net"
	"strconv"
)

type thriftCfg struct {
	localIP   string
	localPort int
	localAddr string
}

func (cfg *thriftCfg) loadConfig() {
	if FileExists(Application.ConfigFile, false) {
		ini, err := inifile.Load(Application.ConfigFile)
		if err != nil {
			Warn("")
			return
		}

		sec := ini.Section("Server")
		vs := sec.Key("localIP").String()
		if vs != "" {
			cfg.localIP = vs
		}

		vi := sec.Key("localPort").MustInt(0)
		if vi > 0 {
			cfg.localPort = vi
		}
		cfg.localAddr = net.JoinHostPort(cfg.localIP, strconv.Itoa(cfg.localPort))
	}
}

type thriftAction struct {
	//thrift handler
}

func (ta *thriftAction) Action(ctx context.Context, param *ActionParam) (_r *ActionResult_, _err error) {
	Info(param.Data)
	return &ActionResult_{
		Res:  true,
		Code: 0,
		Data: "i am server",
	}, nil
}

func (ta *thriftAction) ActionClient(ctx context.Context, param *ActionParam) (_err error) {
	Info(param.Data)
	return nil
}

//--------------------------------------------------------------------------------

//thriftServer thrift服务对象
var thriftServer *thrift.TSimpleServer = nil

/*StartThriftService 2022-06-09 13:39:10
  描述: 启动thrift服务
*/
func StartThriftService() {
	cfg := thriftCfg{
		localIP:   "",
		localPort: 8080,
	}

	cfg.loadConfig()
	//apply config file
	Info("thrift service starting...")

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(cfg.localAddr)
	if err != nil {
		Error("thrift.NewTServerSocket failure: " + err.Error())
		return
	}

	handler := &thriftAction{}
	processor := NewBusinessProcessor(handler)

	thriftServer = thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	Info("thrift service on: " + cfg.localAddr)
	thriftServer.Serve()
}

/*StopThriftService 2022-06-09 13:44:26
  描述: 关闭thrift服务
*/
func StopThriftService() (err error) {
	defer ErrorHandle(false, func(e any) {
		if e != nil {
			err = errors.New("stop thrift service failure.")
		}
	})

	if thriftServer != nil {
		thriftServer.Stop()
		thriftServer = nil
		Info("thrift service closed.")
	}

	return nil
}
