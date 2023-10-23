package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ANamelessSystem/go-mahobotlite/botlogic"
	"github.com/sirupsen/logrus"
)

func main() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logrus.SetLevel(logrus.DebugLevel)
	botlogic.SetCQAddr(CqHttpAddr)
	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile)) // Multi log output: file and stdout(terminal)

	http.HandleFunc("/", botlogic.ReceiveHandler) // 注册处理函数
	info, errInfo := botlogic.GetLoginInfo()
	if errInfo != nil {
		logrus.Errorf("请求发送失败, BOT->CQ无法建立连接: %v", errInfo)
		return
	}
	logrus.Infof("请求发送成功, BOT->CQ连接建立. 登录ID: %v, 名称: %v", info.Data.UserID, info.Data.Nickname)

	go func() {
		for {
			select {
			case <-botlogic.HeartbeatReceived:
				if currentConnectionState == Disconnected {
					logrus.Info("成功收到CQ心跳消息, CQ->BOT连接建立")
					currentConnectionState = Connected
				}
				receivedHeartbeat = true
			case <-time.After(10 * time.Second):
				if receivedHeartbeat && currentConnectionState == Disconnected {
					logrus.Info("收到CQ心跳消息, CQ->BOT连接恢复")
					currentConnectionState = Connected
				} else if !receivedHeartbeat && currentConnectionState == Connected {
					logrus.Warn("未收到CQ心跳消息, CQ->BOT连接中断")
					currentConnectionState = Disconnected
				}
				receivedHeartbeat = false
			}
		}
	}()

	logrus.Infof("启动服务器... 监听端口: %s", LsnrAddr)

	if err := http.ListenAndServe(LsnrAddr, nil); err != nil {
		panic(err)
	}
}

type ConnectionState int

const (
	Disconnected ConnectionState = iota
	Connected
)

var currentConnectionState = Disconnected
var receivedHeartbeat bool
