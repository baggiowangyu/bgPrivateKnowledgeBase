# 使用GoLand基于dubbo-go开发Dubbo服务端

## 0x01 创建项目

<font color="red">**注意在创建项目的时候要选择 Go Modules (vgo)那个**</font>

![](assets/005/20191017-545ad263.png)  

项目生成后，会自动生成 ```go.mod``` 文件，此文件用于管理本项目的依赖库。

![](assets/005/20191017-9d69da15.png)  

参考dubbo-go的范例，在go.mod中添加以下内容：

```
require (
	github.com/Workiva/go-datastructures v1.0.50
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20190802083043-4cd0c391755e // indirect
	github.com/apache/dubbo-go v1.1.0
	github.com/apache/dubbo-go-hessian2 v1.2.5-0.20190909140437-80cbb25cbb22
	github.com/buger/jsonparser v0.0.0-20181115193947-bf1c66bbce23 // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/creasty/defaults v1.3.0
	github.com/dubbogo/getty v1.2.2
	github.com/dubbogo/gost v1.1.1
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/golang/mock v1.3.1
	github.com/google/btree v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.5 // indirect
	github.com/hashicorp/consul v1.5.3
	github.com/hashicorp/consul/api v1.1.0
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jinzhu/copier v0.0.0-20190625015134-976e0346caa8
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/lestrrat/go-envload v0.0.0-20180220120943-6ed08b54a570 // indirect
	github.com/lestrrat/go-file-rotatelogs v0.0.0-20180223000712-d3151e2a480f // indirect
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/nacos-group/nacos-sdk-go v0.0.0-20190723125407-0242d42e3dbb
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/goconvey v0.0.0-20190710185942-9d28bd7c0945 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/tebeka/strftime v0.1.3 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/toolkits/concurrent v0.0.0-20150624120057-a4371d70e3e3 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.etcd.io/etcd v3.3.13+incompatible
	go.uber.org/atomic v1.4.0
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.22.1
	gopkg.in/yaml.v2 v2.2.2
)
```

GoLand会开始下载相关版本的库到%GOPATH%/pkg/mod目录中，若已下载，则不用再重复下载。下载完成后，项目工程的 ```External Libraries``` 节点下```Go Modules``` 节点里面存放的是我们指定的依赖库源码。具体如下图：

![](assets/005/20191017-0805d832.png)  

到此，项目创建已经完毕。接下来开始编写服务代码。

## 0x02 编写Dubbo服务端实现代码

我们这里的源代码目录结构，参考dubbo-go的example内的目录结构：

- app：源代码目录
- assembly：编译/运行脚本目录
- profiles：配置文件目录

### app目录

在这个例子里面，我们要写三个源码文件：

- 框架代码：
- Dubbo服务代码：
- 版本代码：


#### Dubbo服务代码

##### 交互数据结构定义(user.go)

交互数据结构主要是用于在Dubbo服务器与Dubbo客户端之间传递的数据结构。
由于Dubbo原生是Java实现的，当我们与Java交互时需要根据Java特性实现对应的数据类型。

其具体代码如下：

```
package main

import (
	"fmt"
	"github.com/apache/dubbo-go-hessian2"
	"strconv"
	"time"
)

//////////////////////////////////////////////////////////
// 这里定义的三块内容，主要用于Java枚举类型的处理
// 要定义其他的枚举类型的基本操作也是这样
const (
	MAN hessian.JavaEnum = iota
	WOMAN
)

var genderName = map[hessian.JavaEnum]string{
	MAN:	"MAN",
	WOMAN:	"WOMAN",
}

var genderValue = map[string]hessian.JavaEnum{
	"MAN":		MAN,
	"WOMAN":	WOMAN,
}
//////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////
// 这里定义一个基于Java的枚举类型对象Gender，以及实现其应有的接口
type Gender hessian.JavaEnum

func (g Gender) JavaClassName() string {
	return "com.ikurento.user.Gender"
}

// 根据枚举值获取枚举名称
func (g Gender) String() string {
	s, ok := genderName[hessian.JavaEnum(g)]
	if ok {
		return s
	}

	return strconv.Itoa(int(g))
}

// 根据枚举名称获取枚举值
func (g Gender) EnumValue(s string) hessian.JavaEnum {
	v, ok := genderValue[s]
	if ok {
		return v
	}

	return hessian.InvalidJavaEnum
}
//////////////////////////////////////////////////////

//////////////////////////////////////////////////////
// 这里定义实际交互数据使用的对象

// 在Go的面向对象中，struct的成员变量：
// 首字母大写的为公有成员
// 首字母小写的为私有成员
// 这里要设置为公有成员，不要弄错了
type User struct {
	Id		string
	Name	string
	Age		int32
	Time	time.Time
	Sex		Gender	// 需要注意的是，这里要存枚举对象，实际上到了go这边就是字符串了
}

func (u User) JavaClassName() string {
	return "com.ikurento.user.User"
}

func (u User) String() string {
	return fmt.Sprintf(
		"User{Id:%s, Name:%s, Age:%d, Time:%s, Sex:%s}",
		u.Id, u.Name, u.Age, u.Time, u.Sex,
	)
}
//////////////////////////////////////////////////////
```

##### 服务提供者实现代码(user_provider.go)

服务提供者这里需要实现具体的接口以及相关逻辑。
需要注意的是，服务提供者需要实现三个基本接口(V1.1.0是这样定义的，后面的版本看代码已经改成了两个)：

- MethodMapper
- Service
- Version

具体实现代码如下：

```
package main

import (
	"context"
	"fmt"
	"github.com/apache/dubbo-go-hessian2/java_exception"
	"github.com/apache/dubbo-go/config"
	"strconv"
	"time"

	perrors "github.com/pkg/errors"
)

//////////////////////////////////////////////////////
// 这一部分根据实际业务需要进行调整，主要目的是生成交互数据
var userMap = UserProvider{user: make(map[string]User)}
var DefaultUser = User{
	Id: 	"0",
	Name: 	"Alex Stocks",
	Age:	31,
	Sex:	Gender(MAN),
}

func init() {
	// 注册Provider
	config.SetProviderService(new(UserProvider))

	userMap.user["A000"] = DefaultUser
	userMap.user["A001"] = User{Id: "001", Name: "ZhangSheng", Age: 18, Sex: Gender(MAN)}
	userMap.user["A002"] = User{Id: "002", Name: "Lily", Age: 20, Sex: Gender(WOMAN)}
	userMap.user["A003"] = User{Id: "113", Name: "Moorse", Age: 30, Sex: Gender(WOMAN)}

	for k, v := range userMap.user {
		v.Time = time.Now()
		userMap.user[k] = v
	}
}
//////////////////////////////////////////////////////

//////////////////////////////////////////////////////
// 服务提供者定义
type UserProvider struct {
	user map[string]User
}

// V1.1.0必备函数：
// - MethodMapper
// - Service
// - Version
func (u *UserProvider) MethodMapper() map[string]string {
	return map[string]string{"GetUser2": "getUser",}
}

func (u *UserProvider) Service() string {
	return "com.ikurento.user.UserProvider"
}

func (u *UserProvider) Version() string {
	return ""
}

// 私有成员
func (u *UserProvider) getUser(userId string) (*User, error) {
	if user, ok := userMap.user[userId]; ok {
		return &user, nil
	}

	return nil, fmt.Errorf("invalid user id:%s", userId)
}

// 公有成员
func (u *UserProvider) GetUser(ctx context.Context, req []interface{}, rsp *User) error {
	var (
		err  error
		user *User
	)

	println("req:%#v", req)
	user, err = u.getUser(req[0].(string))
	if err == nil {
		*rsp = *user
		println("rsp:%#v", rsp)
	}
	return err
}

func (u *UserProvider) GetUser0(id string, name string) (User, error) {
	var err error

	println("id:%s, name:%s", id, name)
	user, err := u.getUser(id)
	if err != nil {
		return User{}, err
	}
	if user.Name != name {
		return User{}, perrors.New("name is not " + user.Name)
	}
	return *user, err
}

func (u *UserProvider) GetUser2(ctx context.Context, req []interface{}, rsp *User) error {
	var err error

	println("req:%#v", req)
	rsp.Id = strconv.Itoa(int(req[0].(int32)))
	return err
}

func (u *UserProvider) GetUser3() error {
	return nil
}

func (u *UserProvider) GetErr(ctx context.Context, req []interface{}, rsp *User) error {
	return java_exception.NewThrowable("exception")
}

func (u *UserProvider) GetUsers(req []interface{}) ([]interface{}, error) {
	var err error

	println("req:%s", req)
	t := req[0].([]interface{})
	user, err := u.getUser(t[0].(string))
	if err != nil {
		return nil, err
	}
	println("user:%v", user)
	user1, err := u.getUser(t[1].(string))
	if err != nil {
		return nil, err
	}
	println("user1:%v", user1)

	return []interface{}{user, user1}, err
}
//////////////////////////////////////////////////////
```

#### 框架代码(server.go)

我们在app目录下创建server.go文件。此文件可以看成是一个相对固定的启动代码，直接复制粘贴使用即可。具体内容如下：

```
package main

import (
	"fmt"
	"github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	// 这些看上去好像没什么用，暂时留着
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

var survivalTimeout = int(3e9)

// dubbo-go的配置文件是通过环境变量被程序加载的，所以每次启动的时候要先设置当前程序需要的环境变量
// Windows平台：
//
// set CONF_PROVIDER_FILE_PATH="xxx"
// set APP_LOG_CONF_FILE="xxx"
//
// Linux平台：
//
// export CONF_PROVIDER_FILE_PATH="xxx"
// export APP_LOG_CONF_FILE="xxx"
//
func main() {
	// ------for hessian2------
	hessian.RegisterJavaEnum(Gender(MAN))
	hessian.RegisterJavaEnum(Gender(WOMAN))
	hessian.RegisterPOJO(&User{})
	// ------------

	config.Load()

	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			go time.AfterFunc(time.Duration(float64(survivalTimeout)*float64(time.Second)), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// 要么fastFailTimeout时间内执行完毕下面的逻辑然后程序退出，要么执行上面的超时函数程序强行退出
			fmt.Println("provider app exit now...")
			return
		}
	}
}
```

#### 版本代码(version.go)

这个文件我现在也不清楚是干什么用的，先照抄，后面大家自己体会吧。

```
package main

var (
	Version = "2.6.0"
)

```

## 0x03 编写Dubbo服务端配置文件以及日志配置文件

服务配置文件和日志配置文件均放在/profiles/dev目录内。

### server.yml

关于配置文件每个参数的定义，可以参考Dubbo官网相关内容。这里就不详细展开了。参考地址：[http://dubbo.apache.org/zh-cn/docs/user/references/xml/dubbo-service.html](http://dubbo.apache.org/zh-cn/docs/user/references/xml/dubbo-service.html)

```
# dubbo server yaml configure file

# application config
application_config:
  organization : "ikurento.com"
  name : "BDTService"
  module : "dubbogo user-info server"
  version : "0.0.1"
  owner : "ZX"
  environment : "dev"

registries :
  "hangzhouzk":
    # 对应java配置中address属性的zookeeper <dubbo:registry address="zookeeper://127.0.0.1:2181"/>
    protocol: "zookeeper"
    timeout	: "3s"
    address: "127.0.0.1:2181"
    username: ""
    password: ""
  "shanghaizk":
    protocol: "zookeeper"
    timeout	: "3s"
    address: "127.0.0.1:2182"
    username: ""
    password: ""


services:
  "UserProvider":
    # 可以指定多个registry，使用逗号隔开;不指定默认向所有注册中心注册
    registry: "hangzhouzk"
    protocol : "dubbo"
    # 相当于dubbo.xml中的interface
    interface : "com.ikurento.user.UserProvider"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    methods:
      - name: "GetUser"
        retries: 1
        loadbalance: "random"

protocols:
  "dubbo1":
    name: "dubbo"
    #    ip : "127.0.0.1"
    port: 20000


protocol_conf:
  dubbo:
    session_number: 700
    fail_fast_timeout: "5s"
    session_timeout: "20s"
    getty_session_param:
      compress_encoding: false
      tcp_no_delay: true
      tcp_keep_alive: true
      keep_alive_period: "120s"
      tcp_r_buf_size: 262144
      tcp_w_buf_size: 65536
      pkg_rq_size: 1024
      pkg_wq_size: 512
      tcp_read_timeout: "1s"
      tcp_write_timeout: "5s"
      wait_timeout: "1s"
      max_msg_len: 1024
      session_name: "server"

```

### log.yml

```
level: "debug"
development: true
disableCaller: false
disableStacktrace: false
sampling:
encoding: "console"

# encoder
encoderConfig:
  messageKey: "message"
  levelKey: "level"
  timeKey: "time"
  nameKey: "logger"
  callerKey: "caller"
  stacktraceKey: "stacktrace"
  lineEnding: ""
  levelEncoder: "capitalColor"
  timeEncoder: "iso8601"
  durationEncoder: "seconds"
  callerEncoder: "short"
  nameEncoder: ""

outputPaths:
  - "stderr"
errorOutputPaths:
  - "stderr"
initialFields:
```

## 0x04 编译与运行

### 服务编译

在app目录内，执行go build，等待app.exe生成；

### 服务启动

由于我是在Windows下开发的，在assembly/dev目录内创建了startup.bat。内容为：

```
set CONF_PROVIDER_FILE_PATH=E:\opensource\DubboGoServer01\profiles\dev\server.yml
set APP_LOG_CONF_FILE=E:\opensource\DubboGoServer01\profiles\dev\log.yml

E:\opensource\DubboGoServer01\app\app.exe
```

接下来在本地准备好zookeeper服务，并启动。

启动startup.bat，即可看到Dubbo服务在zookeeper上注册好服务了。如下图：

![](assets/005/20191017-cac45d3c.png)  

![](assets/005/20191017-f0a61547.png)  
