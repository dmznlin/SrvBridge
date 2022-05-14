package srvbridge

import (
	. "SrvBridge/src/mit"
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	. "github.com/dmznlin/znlib-go/znlib"
	inifile "github.com/go-ini/ini"
	"net"
	"os"
	"strconv"
)

/******************************************************************************
作者: dmzn@163.com 2022-05-14
描述: service for thrift
******************************************************************************/

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
		Data: "",
	}, nil
}

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
		os.Exit(1)
	}

	handler := &thriftAction{}
	processor := NewBusinessProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	Info("thrift service on: " + cfg.localAddr)
	server.Serve()
}
