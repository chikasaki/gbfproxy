package main

import (
	"fmt"
	"github.com/armon/go-socks5"
	"myproxy/comm"
	"myproxy/crypt"
	mylog "myproxy/mylog"
	"strconv"
	"strings"
	"time"
)

type Socks5ProxyConf struct {
	Auth struct {
		NeedAuth           bool
		Username, Password string
		CryptConf          crypt.Config
	}
	Network struct {
		Addr string
	}
}

//	type CredentialStore interface {
//		Valid(user, password string) bool
//	}
type CryptCredentialStore struct {
	username, passwd string
	c                crypt.Crypt
}

type ResultStr struct {
	Str       string
	Timestamp int64
}

func (s *CryptCredentialStore) getStrSplit(str string) (ans ResultStr, err error) {
	splits := strings.Split(str, "_")
	if len(splits) < 2 {
		return ResultStr{}, fmt.Errorf("[Err] invalid input str:%s", str)
	}
	ans.Timestamp, err = strconv.ParseInt(splits[len(splits)-1], 10, 64)
	if err != nil {
		return ResultStr{}, err
	}
	ans.Str = strings.Join(splits[:len(splits)-1], "_")
	return ans, nil
}

func (s *CryptCredentialStore) Valid(encryptUser, encryptPassword string) bool {
	fmt.Println("start to validate")
	var (
		err                      error
		lNow                     = time.Now().Unix()
		decryptUserSequence      = string(s.c.Decrypt([]byte(encryptUser)))     //realUsername_timestamp
		decryptPasswordSequence  = string(s.c.Decrypt([]byte(encryptPassword))) //realPassword_timestamp
		userResult, passwdResult ResultStr
	)

	userResult, err = s.getStrSplit(decryptUserSequence)
	if err != nil {
		mylog.Errorf("[Err] getStrSplit of user fail, input:%s, err:%+v\n", decryptUserSequence, err)
		return false
	}
	passwdResult, err = s.getStrSplit(decryptPasswordSequence)
	if err != nil {
		mylog.Errorf("[Err] getStrSplit of passwd fail, input:%s, err:%+v\n", decryptPasswordSequence, err)
		return false
	}

	//judge time
	if comm.Abs[int64](lNow, userResult.Timestamp) > 3 {
		mylog.Errorf("[Err] server timestamp:%d, request timestamp:%d, too different", userResult.Timestamp, lNow)
		return false
	}
	if userResult.Str != s.username || passwdResult.Str != s.passwd {
		return false
	}
	return true
}

func NewSocks5Proxy(conf *Socks5ProxyConf) {
	var (
		credential socks5.CredentialStore
		socksConf  socks5.Config
		socksSvr   *socks5.Server
		err        error
	)
	if conf.Auth.NeedAuth {
		c, err := crypt.NewCrypt(&conf.Auth.CryptConf)
		if err != nil {
			mylog.Fatalf("[Fatal] new crypt fail, err:%+v", err)
			return
		}
		credential = &CryptCredentialStore{
			username: conf.Auth.Username,
			passwd:   conf.Auth.Password,
			c:        c,
		}
	}
	socksConf = socks5.Config{
		Credentials: credential,
		Rules: &socks5.PermitCommand{
			EnableConnect:   true,
			EnableBind:      false,
			EnableAssociate: false,
		},
		Logger: mylog.Logger,
	}
	socksSvr, err = socks5.New(&socksConf)
	if err != nil {
		mylog.Fatalf("[Fatal] new socks svr fail, err:%+v", err)
		return
	}
	socksSvr.ListenAndServe("tcp", conf.Network.Addr)
}
