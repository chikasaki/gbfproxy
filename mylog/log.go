package mylog

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	LogFile      string
	HoldFilesNum int
}

var (
	Logger     *log.Logger
	nowDateStr string
)

func InitLog(conf Config) {
	if conf.LogFile == "" {
		Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
		return
	}
	go logRotate(conf)
	for Logger == nil {
	}
}

func logRotate(conf Config) {
	for {
		time.Sleep(time.Second)
		if nowDateStr != time.Now().Format("20060102") {
			nowDateStr = time.Now().Format("20060102")
			filename := conf.LogFile + "_" + nowDateStr + ".log"
			file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatal("fail to create log file! err:", err)
				return
			}
			l := log.New(file, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
			Logger = l
		}
	}
}

func Errorf(format string, v ...any) {
	Logger.Output(2, fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...any) {
	Logger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
