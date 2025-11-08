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
	"strings"
	"sync"
	"time"
)

var logger *slog.Logger
var listenAddr *net.TCPAddr
var wg *sync.WaitGroup
var running bool = false
var v1Msgs map[string][]string

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
	v1Msgs = makeV1Messages(config)
	listenAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Tcp.IP, config.Tcp.Port))
	if err != nil {
		logger.Error("Failed to setup TCP", "error", err)
		return false
	}
	return true
}

func makeV1Messages(config *config.Config) map[string][]string {
	result := make(map[string][]string)
	if _, ok := config.Dmx.Groups["stg"]; ok {
		result["fi"] = []string{"fadeIn", "stg"}
		result["fo"] = []string{"fadeOut", "stg"}
		result["ci"] = []string{"fadeIn", "stg", "interval:0"}
		result["co"] = []string{"fadeOut", "stg", "interval:0"}
	}
	if _, ok := config.Dmx.Groups["aud"]; ok {
		result["fai"] = []string{"fadeIn", "aud"}
		result["fao"] = []string{"fadeOut", "aud"}
	}
	result["mute"] = []string{"mute", "true"}
	result["unmute"] = []string{"mute", "false"}
	return result
}

func handleRequest(conn *net.TCPConn) {
	logger.Info("Connect", "remote", conn.RemoteAddr())
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		l, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Warn("Invalid close", "error", err)
			}
			break
		}
		msgs := regexp.MustCompile("\r\n|\n|\r").Split(string(buf[:l]), -1)
		for _, msg := range msgs {
			cmd := strings.Split(msg, " ")
			if v1Msgs != nil && len(cmd) == 1 {
				for key, newCmd := range v1Msgs {
					if cmd[0] == key {
						cmd = newCmd
						break
					}
				}
			}
			logger.Debug("TCP message", "cmd", cmd)
			switch cmd[0] {
			case "fadeIn", "fadeOut":
				isIn := cmd[0] == "fadeIn"
				msgArg := message.MessageBody{
					Action: "fade",
					Arg:    map[string]string{},
				}
				if len(cmd) <= 1 {
					continue
				}
				if len(cmd) >= 3 {
					args := strings.Split(cmd[2], ",")
					for _, v := range args {
						if !strings.Contains(v, ":") {
							continue
						}
						arg := strings.Split(v, ":")
						logger.Debug("test", "arg", arg)
						msgArg.Arg[arg[0]] = arg[1]
					}
				}
				msgArg.Arg["id"] = cmd[1]
				msgArg.Arg["isIn"] = fmt.Sprintf("%v", isIn)
				go packageModule.ModuleManager.SendMessage(message.Message{
					To:  "dmx",
					Arg: msgArg,
				})
			case "mute":
				mute := true
				if len(cmd) >= 2 && cmd[1] == "false" {
					mute = false
				}
				go packageModule.ModuleManager.SendMessage(message.Message{
					To: "osc",
					Arg: message.MessageBody{
						Action: "mute",
						Arg: map[string]string{
							"mute": fmt.Sprintf("%v", mute),
						},
					},
				})
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
