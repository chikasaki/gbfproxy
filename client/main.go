package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"myproxy/mylog"
)

var gConf = struct {
	LogConf mylog.Config
	Socks5  struct {
		ProxyConf Socks5ProxyConf
		PoolConf  Sock5ConnPoolConf
	}
}{}

func main() {
	_, err := toml.DecodeFile("client_config.toml", &gConf)
	if err != nil {
		log.Fatal("decode config file error:", err)
	}

	mylog.InitLog(gConf.LogConf)
	NewSocks5Proxy(gConf.Socks5.ProxyConf, gConf.Socks5.PoolConf)
}
