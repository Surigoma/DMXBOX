package artnet

import (
	"backend/config"
	"bytes"
	"fmt"
	"log/slog"
	"net"
	"time"
)

type StaticBase struct {
	srcPort uint
	rmtPort uint
	Version []uint8
	Ops     map[string][]byte
}

var static StaticBase = StaticBase{
	srcPort: 6454,
	rmtPort: 6454,
	Version: []uint8{
		0,  // High
		14, // Low
	},
	Ops: map[string][]byte{
		"OpDMX":     {0x00, 0x50},
		"OpPoll":    {0x00, 0x20},
		"OpPollRep": {0x00, 0x21},
	},
}

type ArtnetAddress struct {
	Net         uint8
	SubUniverse uint8
	Universe    uint8
}
type Artnet struct {
	TargetAddr string
	targetUDP  *net.UDPAddr
	sourceUDP  *net.UDPAddr
	socket     *net.UDPConn
	logger     *slog.Logger
	isRunning  bool
	Running    bool
	address    ArtnetAddress
}

func (a *Artnet) Initialize(log *slog.Logger, config *config.Config) bool {
	var err error
	a.logger = log
	a.targetUDP, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", a.TargetAddr, static.rmtPort))
	if err != nil {
		a.logger.Error("Failed initialize Artnet System.", "err", err)
		return false
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		a.logger.Error("Failed to get Interface ips.", "err", err)
		return false
	}
	for _, addr := range addrs {
		_, cidr, _ := net.ParseCIDR(addr.String())
		if cidr.Contains(a.targetUDP.IP) {
			a.sourceUDP, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr.(*net.IPNet).IP.String(), static.srcPort))
			if err != nil {
				a.logger.Error("failed to create", "addr", addr)
			}
			break
		}
	}
	if a.sourceUDP == nil {
		a.logger.Error("Failed to contain network for interfaces.", "target", a.targetUDP)
		return false
	}
	a.address = ArtnetAddress{
		Net:         config.Output.Artnet.Net,
		SubUniverse: config.Output.Artnet.SubUniverse,
		Universe:    config.Output.Artnet.Universe,
	}
	a.logger.Debug("Initialize artnet", "src", a.sourceUDP, "dst", a.targetUDP, "address", a.address)
	return true
}

func (a *Artnet) listen() {
	var buffer []byte = make([]byte, 1024)
	a.Running = true
	for a.isRunning {
		n, err := a.socket.Read(buffer)
		if err != nil {
			a.isRunning = false
			break
		}
		if !bytes.Equal(buffer[0:7], []byte("Art-Net")) {
			continue
		}
		op := buffer[8:10]
		cmd := ""
		for k, v := range static.Ops {
			if bytes.Equal(op, v) {
				if k != "OpDMX" {
					a.logger.With("dir", "recv").Debug("CMD: "+k, "n", n, "data", buffer[:n])
				}
				cmd = k
				break
			}
		}
		if cmd == "" {
			a.logger.With("dir", "recv").Debug("OpUNKNOWN", "n", n, "data", buffer[:n])
		}
	}
	a.Running = false
}

func (a *Artnet) Start() bool {
	var err error
	a.socket, err = net.ListenUDP("udp", a.sourceUDP)
	if err != nil {
		a.logger.Error("Failed to create socket", "err", err, "src", a.sourceUDP, "dst", a.targetUDP)
		return false
	}
	a.logger.Debug("Start")
	a.isRunning = true
	go a.listen()
	return true
}

func (a *Artnet) Stop() bool {
	a.isRunning = false
	if a.socket != nil {
		a.socket.Close()
	}
	for range 100 {
		if !a.Running {
			break
		}
		time.After(10 * time.Millisecond)
	}
	a.logger.Debug("Stop")
	return true
}

func (a *Artnet) RenderData(op string, data *[]byte) *[]byte {
	result := []byte("Art-Net")
	result = append(result, 0x0)
	result = append(result, static.Ops[op]...)
	result = append(result, static.Version...)
	result = append(result, *data...)
	return &result
}

var sequence uint8 = 0x0

func (a *Artnet) SendDMXData(data *[]byte) {
	length := uint16(len(*data))
	before := []byte{sequence}
	before = append(before, 0, ((a.address.SubUniverse&0xf)<<4)|(a.address.Universe&0xf), a.address.Net)
	before = append(before, byte((length>>8)&0xff), byte(length&0xff))
	before = append(before, *data...)
	renderd := a.RenderData("OpDMX", &before)
	sequence += 1
	if sequence == 0 {
		sequence = 1
	}
	go func() {
		_, err := a.socket.WriteToUDP(*renderd, a.targetUDP)
		if err != nil {
			a.logger.Error("Drop", "data", renderd)
			return
		}
	}()
}
