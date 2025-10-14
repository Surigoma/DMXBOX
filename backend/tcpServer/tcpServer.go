package tcpserver

import (
	"backend/message"
	"backend/packageModule"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"
)

type TCPParams struct {
	logger     slog.Logger
	listenAddr *net.TCPAddr
}

var p TCPParams
var channel chan message.Message
var wg *sync.WaitGroup
var running bool = true

var TcpServer packageModule.PackageModule = packageModule.PackageModule{
	Initialize: Initialize,
	Run:        StartTCP,
}

func Initialize(param packageModule.PackageModuleParam) bool {
	var err error
	p.logger = param.Logger
	wg = param.Wg
	p.listenAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", param.Config.Tcp.IP, param.Config.Tcp.Port))
	if err != nil {
		p.logger.Error("Failed to setup TCP", "error", err)
		return false
	}
	return true
}
func handleRequest(conn *net.TCPConn) {
	p.logger.Info("Connect", "remote", conn.RemoteAddr().Network())
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				p.logger.Warn("Invalid close", "error", err)
			}
			break
		}
		msgs := strings.Split(string(buf), "\n")
		for _, msg := range msgs {
			switch msg {
			case "test":
				p.logger.Info("test")
			}
			conn.Write([]byte("ack\n"))
		}
	}
}
func tcpThread(ln *net.TCPListener) {
	defer wg.Done()
	for running {
		err := ln.SetDeadline(time.Now().Add(time.Second * 10))
		if err != nil {
			p.logger.Error("Failed to set Dead line.", "error", err)
			break
		}
		conn, err := ln.AcceptTCP()
		if err != nil {
			p.logger.Error("Failed setup connection", "error", err)
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

func messageProcess() {
	for {
		mes := <-channel
		if mes.To == "tcp" {
			if handleMessage(mes) < 0 {
				break
			}
		}
	}
}

func StartTCP() {
	p.logger.Info("Hello TCP server", "listenAddr", p.listenAddr)
	ln, err := net.ListenTCP("tcp", p.listenAddr)
	if err != nil {
		p.logger.Error("Failed to start a tcp server", "error", err)
		return
	}
	wg.Add(1)
	defer ln.Close()
	go tcpThread(ln)
	go messageProcess()
}
