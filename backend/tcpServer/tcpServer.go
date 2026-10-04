package tcpserver

import (
	"backend/config"
	"backend/message"
	"backend/packageModule"
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"
)

var logger *slog.Logger
var listenAddr *net.TCPAddr
var wg *sync.WaitGroup
var runningWg sync.WaitGroup
var running bool = false
var runningMutex sync.Mutex
var listener *net.TCPListener
var listenerMutex sync.Mutex
var v1Msgs map[string][]string
var currentModule *packageModule.PackageModule
var connections = make(map[*net.TCPConn]struct{})
var connectionsMutex sync.Mutex
var connectionWg sync.WaitGroup

var TcpServer packageModule.PackageModule = packageModule.PackageModule{
	ModuleName:     "tcp",
	Initialize:     Initialize,
	Run:            StartTCP,
	Stop:           StopTCP,
	MessageHandler: HandleMessage,
}

func changeRunning(run bool) {
	runningMutex.Lock()
	running = run
	runningMutex.Unlock()
}

func isRunning() bool {
	runningMutex.Lock()
	defer runningMutex.Unlock()
	return running
}

func changeListener(ln *net.TCPListener) {
	listenerMutex.Lock()
	listener = ln
	listenerMutex.Unlock()
}

func closeListener() {
	listenerMutex.Lock()
	defer listenerMutex.Unlock()
	if listener != nil {
		_ = listener.Close()
		listener = nil
	}
}

func Initialize(module *packageModule.PackageModule, config *config.Config) bool {
	var err error
	currentModule = module
	logger = module.Logger
	wg = module.Wg
	runningWg = sync.WaitGroup{}
	runningMutex = sync.Mutex{}
	listenerMutex = sync.Mutex{}
	changeListener(nil)
	changeRunning(true)
	v1Msgs = makeV1Messages(config)
	listenAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Input.Tcp.IP, config.Input.Tcp.Port))
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
	manager := packageModule.GetModuleManager()
	logger.Info("Connect", "remote", conn.RemoteAddr())
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 512), 4098)
	scanner.Split(splitCommand)
	for scanner.Scan() {
		cmd := strings.Fields(scanner.Text())
		if len(cmd) == 0 {
			continue
		}
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
			currentModule.SendMessage(message.Message{
				To:  "dmx",
				Arg: msgArg,
			})
		case "mute":
			mute := true
			if len(cmd) >= 2 && cmd[1] == "false" {
				mute = false
			}
			currentModule.SendMessage(message.Message{
				To: "osc",
				Arg: message.MessageBody{
					Action: "mute",
					Arg: map[string]string{
						"isMute": fmt.Sprintf("%v", mute),
					},
				},
			})
		case "test":
			logger.Debug("test")
			manager.SendMessage(message.Message{
				To: "test",
			})
		}
		conn.Write([]byte("ack\r\n"))
	}
	if err := scanner.Err(); err != nil {
		logger.Warn("Invalid TCP command", "error", err)
	}
	logger.Info("Disconnect", "remote", conn.RemoteAddr())
}

// CRLF is one separator even if CR and LF arrive in different reads.
func splitCommand(data []byte, atEOF bool) (int, []byte, error) {
	for i, b := range data {
		if i > 4095 {
			return 0, nil, errors.New("TCP command exceeds 4095 bytes")
		}
		if b == '\n' {
			return i + 1, data[:i], nil
		}
		if b == '\r' {
			advance := i + 1
			if advance < len(data) && data[advance] == '\n' {
				advance++
			}
			return advance, data[:i], nil
		}
	}
	if atEOF && len(data) > 0 {
		if len(data) > 4095 {
			return 0, nil, errors.New("TCP command exceeds 4095 bytes")
		}
		return len(data), data, nil
	}
	return 0, nil, nil
}
func tcpThread(ln *net.TCPListener) {
	defer wg.Done()
	defer runningWg.Done()
	defer closeListener()
	for isRunning() {
		conn, err := ln.AcceptTCP()
		if err != nil {
			if !isRunning() {
				break
			}
			logger.Error("Failed setup connection", "error", err)
			continue
		}
		connectionsMutex.Lock()
		connections[conn] = struct{}{}
		connectionWg.Add(1)
		connectionsMutex.Unlock()
		go func() {
			defer connectionWg.Done()
			defer func() { connectionsMutex.Lock(); delete(connections, conn); connectionsMutex.Unlock() }()
			handleRequest(conn)
		}()
	}
	logger.Info("Close TCP Server")
}

func HandleMessage(mes message.Message) int {
	switch mes.Arg.Action {
	case "reload":
		return 1
	case "stop":
		return -1
	}
	return 0
}

func StartTCP() {
	logger.Info("Hello TCP server", "listenAddr", listenAddr)
	ln, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		logger.Error("Failed to start a tcp server", "error", err)
		changeRunning(false)
		wg.Done()
		return
	}
	changeListener(ln)
	runningWg.Add(1)
	go tcpThread(ln)
}

func StopTCP() {
	if isRunning() {
		changeRunning(false)
		closeListener()
		runningWg.Wait()
		connectionsMutex.Lock()
		for conn := range connections {
			conn.Close()
		}
		connectionsMutex.Unlock()
		connectionWg.Wait()
	}
}
