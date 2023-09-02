package main

import (
	"github.com/armon/go-socks5"
	"myproxy/mylog"
)

type Socks5ProxyConf struct {
	Network struct {
		Addr string
	}
}

func NewSocks5Proxy(conf Socks5ProxyConf, poolConf Sock5ConnPoolConf) {
	var (
		socksConf socks5.Config
		socksSvr  *socks5.Server
		socksPool *Sock5ConnPool
		err       error
	)

	socksPool = NewSock5ConnPool(poolConf)
	socksConf = socks5.Config{
		Rules: &socks5.PermitCommand{
			EnableConnect:   true,
			EnableBind:      false,
			EnableAssociate: false,
		},
		Logger: mylog.Logger,
		Dial:   socksPool.Dial,
	}
	socksSvr, err = socks5.New(&socksConf)
	if err != nil {
		mylog.Fatalf("[Fatal] new socks svr fail, err:%+v", err)
		return
	}
	socksSvr.ListenAndServe("tcp", conf.Network.Addr)
}
