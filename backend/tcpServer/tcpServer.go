package tcpserver

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	"fmt"
	"io"
	"log/slog"
	"net"
	"regexp"
	"sync"
	"time"
)

var logger *slog.Logger
var listenAddr *net.TCPAddr
var wg *sync.WaitGroup
var running bool = false

var TcpServer packageModule.PackageModule = packageModule.PackageModule{
	ModuleName:     "tcp",
	Initialize:     Initialize,
	Run:            StartTCP,
	MessageHandler: handleMessage,
}

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	var err error
	logger = module.Logger
	wg = module.Wg
	running = true
	listenAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Tcp.IP, config.Tcp.Port))
	if err != nil {
		logger.Error("Failed to setup TCP", "error", err)
		return false
	}
	return true
}
func handleRequest(conn *net.TCPConn) {
	logger.Info("Connect", "remote", conn.RemoteAddr())
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		len, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Warn("Invalid close", "error", err)
			}
			break
		}
		msgs := regexp.MustCompile("\r\n|\n|\r").Split(string(buf[:len]), -1)
		fmt.Println(msgs)
		for _, msg := range msgs {
			switch msg {
			case "test":
				logger.Debug("test")
				go packageModule.ModuleManager.SendMessage(message.Message{
					To: "test",
				})
			}
			conn.Write([]byte("ack\r\n"))
		}
	}
	logger.Info("Disconnect", "remote", conn.RemoteAddr())
}
func tcpThread(ln *net.TCPListener) {
	defer wg.Done()
	defer ln.Close()
	for running {
		err := ln.SetDeadline(time.Now().Add(time.Second))
		if err != nil {
			logger.Error("Failed to set Dead line.", "error", err)
			break
		}
		conn, err := ln.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); !ok {
				logger.Error("Failed setup connection", "error", opErr)
			}
			continue
		}
		go handleRequest(conn)
	}
}

func handleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "stop":
		running = false
		return -1
	}
	return 0
}

func StartTCP() {
	logger.Info("Hello TCP server", "listenAddr", listenAddr)
	ln, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		logger.Error("Failed to start a tcp server", "error", err)
		return
	}
	go tcpThread(ln)
}
