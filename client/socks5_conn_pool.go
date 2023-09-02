package main

import (
	"context"
	"fmt"
	"github.com/txthinking/socks5"
	"io"
	"myproxy/crypt"
	"myproxy/mylog"
	"myproxy/mypool"
	"net"
	"time"
)

type Sock5ConnPoolConf struct {
	TargetSocks5Conf struct {
		Addr string
	}
	Auth struct {
		NeedAuth           bool
		Username, Password string
		CryptConf          crypt.Config
	}
	Pool struct {
		ConnCount int
	}
}

type Sock5ConnPool struct {
	pool *mypool.OneTimePool
}

func NewSock5ConnPool(conf Sock5ConnPoolConf) *Sock5ConnPool {
	c, err := crypt.NewCrypt(&conf.Auth.CryptConf)
	if err != nil {
		mylog.Fatalf("new crypt failed, err:%+v", err)
	}
	//预先和server socks5Proxy建立一批连接
	pool := mypool.NewOneTimePool(conf.Pool.ConnCount, func() interface{} {
		ans, err := newSock5Conn(conf, c)
		if err == nil {
			return ans
		}
		return nil
	})
	return &Sock5ConnPool{
		pool: pool,
	}
}

func (p *Sock5ConnPool) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	c := p.pool.Take()
	if c == nil {
		return nil, fmt.Errorf("[ERR] can't get conn")
	}
	client := c.(*socks5.Client)
	if isConnClose(client) {
		//try again
		c = p.pool.Take()
		if c == nil {
			return nil, fmt.Errorf("[ERR] can't get conn")
		}
		client = c.(*socks5.Client)
		if isConnClose(client) {
			return nil, fmt.Errorf("[ERR] can't get a normal tcp connection")
		}
	}
	addrType, parseAddr, parsePort, err := socks5.ParseAddress(addr)
	if err != nil {
		return nil, err
	}
	if addrType == socks5.ATYPDomain {
		addr = addr[1:]
	}
	if _, err = client.Request(socks5.NewRequest(socks5.CmdConnect, addrType, parseAddr, parsePort)); err != nil {
		return nil, err
	}
	return client.TCPConn, nil
}

func newSock5Conn(conf Sock5ConnPoolConf, c crypt.Crypt) (*socks5.Client, error) {
	var (
		lNow                       = time.Now().Unix()
		encryptUser, encryptPasswd string
		client                     *socks5.Client
		err                        error
	)
	if conf.Auth.NeedAuth {
		encryptUser = string(c.Encrypt([]byte(fmt.Sprintf("%s_%d", conf.Auth.Username, lNow))))
		encryptPasswd = string(c.Encrypt([]byte(fmt.Sprintf("%s_%d", conf.Auth.Password, lNow))))
	}
	client, _ = socks5.NewClient(
		conf.TargetSocks5Conf.Addr,
		encryptUser,
		encryptPasswd,
		0, 0)
	// establish connect
	err = client.Negotiate(nil)
	if err != nil {
		mylog.Errorf("[Err] establish sock5 connection failed, err:%+v\n", err)
		return nil, err
	}
	tcpConn := client.TCPConn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(time.Second * 3)
	return client, nil
}

func isConnClose(c *socks5.Client) bool {
	tcpConn, _ := c.TCPConn.(*net.TCPConn)
	if _, err := tcpConn.Read([]byte{}); err != nil {
		if err == io.EOF {
			mylog.Errorf("[Err] tcpConn close")
		}
		mylog.Errorf("[Err] tcpConn read error:%+v", err)
		return true
	}
	return false
}
