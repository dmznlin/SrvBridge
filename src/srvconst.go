/*Package srvbridge ***********************************************************
  作者: dmzn@163.com 2022-06-09 11:55:04
  描述: 公用常量变量定义
******************************************************************************/
package srvbridge

import (
	"context"
	"encoding/json"
	. "github.com/dmznlin/znlib-go/znlib"
	"github.com/go-ini/ini"
	"math/rand"
	"strconv"
)

//ServiceType 支持的桥接服务类型
type ServiceType = uint8

const (
	Srv_redis     ServiceType = iota //redis
	Srv_kafka                        //kafka
	Srv_zookeeper                    //zookeeper
)

//ServiceWorker 服务提供者
type ServiceWorker interface {
	WorkName() string     //服务名称
	LoadConfig(*ini.File) //加载配置
	Start()               //启动函数
	Stop() error          //停止函数
}

var (
	ServerName     string             //服务实例名称
	ServiceWorkers []ServiceWorker    //服务对象列表
	ServiceContext context.Context    //服务上下文
	ServiceCancel  context.CancelFunc //取消上下文
)

//初始化znlib库,且在各单元init()前初始化全局变量
var _ = InitLib(nil, func() {
	ServerName = "ServiceBridge" + strconv.Itoa(rand.Int())
	ServiceWorkers = make([]ServiceWorker, 0, 5)
	ServiceContext, ServiceCancel = context.WithCancel(context.Background())
})

/*RegistWorker 2022-06-10 10:33:08
  参数: sw,服务对象
  描述: 注册一个服务对象
*/
func RegistWorker(sw ServiceWorker) {
	Application.SyncLock.Lock()
	defer Application.SyncLock.Unlock()
	ServiceWorkers = append(ServiceWorkers, sw)
}

/*LoadWorkersConfig 2022-08-17 16:13:02
  描述: 载入配置信息
*/
func LoadWorkersConfig() {
	if FileExists(Application.ConfigFile, false) == false {
		return
	}

	cfg, err := ini.Load(Application.ConfigFile)
	if err != nil {
		Warn(ErrorMsg(err, "SrvBridge.LoadWorkersConfig"))
		return
	}

	str := StrTrim(cfg.Section("config").Key("ServerName").String())
	if str != "" {
		ServerName = str
	}

	for _, worker := range ServiceWorkers {
		worker.LoadConfig(cfg)
	}
}

//--------------------------------------------------------------------------------

const (
	TrancertTag    = ";;;" //
	TrancertTagLen = len(TrancertTag)
)

type ServiceData struct {
	SrvType  ServiceType `json:"SrvType" xml:"SrvType"`   //服务类型
	SrvData  string      `json:"SrvData" xml:"SrvData"`   //数据
	Trancert string      `json:"Trancert" xml:"Trancert"` //跟踪日志
}

/*AddTrancert 2022-06-17 15:26:46
  参数: str,日志数据
  描述: 添加跟踪日志
*/
func (sd *ServiceData) AddTrancert(str string) {
	idx := StrPos(sd.Trancert, ServerName)
	if idx >= 0 {
		tag := StrPos(sd.Trancert, TrancertTag)
		if tag >= 0 {
			sd.Trancert = StrDel(sd.Trancert, idx, tag+TrancertTagLen-1)
		}
	}

	if sd.Trancert == "" {
		sd.Trancert = ServerName + "=" + str + TrancertTag
	} else {
		sd.Trancert = sd.Trancert + ServerName + "=" + str + TrancertTag
	}
}

/*Serialize 2022-06-20 14:07:24
  描述: 序列化结构sd为字符串
*/
func (sd *ServiceData) Serialize() (str string) {
	defer DeferHandle(false, "ServiceData", func(err any) {
		if err != nil {
			str = ""
		}
	})

	data, err := json.MarshalIndent(sd, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(data)
}

/*Deserialize 2022-06-20 14:07:56
  对象: str,数据
  描述: 反序列化str为结构数据
*/
func (sd *ServiceData) Deserialize(str string) (res bool) {
	defer DeferHandle(false, "ServiceData", func(err any) {
		if err != nil {
			res = false
		}
	})

	err := json.Unmarshal([]byte(str), sd)
	if err != nil {
		panic(err)
	}
	return true
}
