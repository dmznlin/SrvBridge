/*Package srvbridge ***********************************************************
  作者: dmzn@163.com 2022-06-09 13:03:31
  描述: service for thrift
******************************************************************************/
package srvbridge

import (
	. "SrvBridge/src/mit"
	"context"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	. "github.com/dmznlin/znlib-go/znlib"
	"github.com/go-ini/ini"
	"net"
	"strconv"
)

//thriftService thrift服务
type thriftService struct {
	enable  bool
	srvIP   string
	srvPort int
	srvAddr string
}

var (
	//thriftWorker thrift服务对象
	thriftWorker *thriftService
	//thriftServer thrift服务器对象
	thriftServer *thrift.TSimpleServer = nil
)

//注册thrift服务
func init() {
	thriftWorker = new(thriftService)
	RegistWorker(thriftWorker)
}

func (srv thriftService) WorkName() string {
	return "thrift-service"
}

func (srv *thriftService) LoadConfig(cfg *ini.File) {
	if len(Application.HostIP) > 0 {
		srv.srvIP = Application.HostIP[0]
	} else {
		srv.srvIP = ""
	}

	srv.enable = true
	srv.srvPort = 8080

	if cfg == nil && FileExists(Application.ConfigFile, false) {
		var err error
		cfg, err = ini.Load(Application.ConfigFile)
		if err != nil {
			Warn(ErrorMsg(err, "thrift.loadconfig"))
			return
		}
	}

	if cfg == nil {
		Warn("thrift.loadconfig: invalid config file.")
		return
	}

	sec := cfg.Section("thrift")
	srv.enable = sec.Key("enable").In("true", []string{"true", "false"}) == "true"
	//是否启动服务

	vs := StrTrim(sec.Key("srvIP").String())
	if vs != "" {
		srv.srvIP = vs
	}

	vi := sec.Key("srvPort").MustInt(0)
	if vi > 0 {
		srv.srvPort = vi
	}
	srv.srvAddr = net.JoinHostPort(srv.srvIP, strconv.Itoa(srv.srvPort))
}

/*Start 2022-06-09 13:39:10
  描述: 启动thrift服务
*/
func (srv thriftService) Start() {
	if !thriftWorker.enable {
		return //不启动服务
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(thriftWorker.srvAddr)
	if err != nil {
		Error("thrift.NewTServerSocket failure: " + err.Error())
		return
	}

	handler := &thriftAction{}
	processor := NewBusinessProcessor(handler)

	thriftServer = thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	Info("thrift service on: " + thriftWorker.srvAddr)
	thriftServer.Serve()
}

/*Stop 2022-06-09 13:44:26
  描述: 关闭thrift服务
*/
func (srv thriftService) Stop() (err error) {
	defer DeferHandle(false, "thriftService.Stop", func(e any) {
		if e != nil {
			err = errors.New(fmt.Sprintf("stop %s failure.", srv.WorkName()))
		}
	})

	if thriftServer != nil {
		thriftServer.Stop()
		thriftServer = nil
		Info(srv.WorkName() + " closed.")
	}

	return nil
}

//--------------------------------------------------------------------------------

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
