package oscserver_test

import (
	"backend/config"
	"backend/message"
	oscserver "backend/oscServer"
	"backend/packageModule"
	"log/slog"
	"net"
	"testing"
	"time"
)

func TestMuteState(t *testing.T) {
	listener, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	module := &packageModule.PackageModule{Logger: slog.Default()}
	cfg := &config.Config{Output: config.OutputTargets{Osc: config.OSCServer{
		Ip: "127.0.0.1", Port: uint16(listener.LocalAddr().(*net.UDPAddr).Port), Type: "int", Format: "/mute/{}", Channels: []uint{1},
	}}}
	if !oscserver.Initialize(module, cfg) {
		t.Fatal("initialize failed")
	}
	if oscserver.GetMuteState().IsMute != nil {
		t.Fatal("initial state must be unknown")
	}
	for _, muted := range []bool{true, false} {
		arg := "false"
		if muted {
			arg = "true"
		}
		oscserver.HandleMessage(message.Message{Arg: message.MessageBody{Action: "mute", Arg: map[string]string{"isMute": arg}}})
		listener.SetReadDeadline(time.Now().Add(time.Second))
		if _, _, err := listener.ReadFrom(make([]byte, 1024)); err != nil {
			t.Fatal(err)
		}
		state := oscserver.GetMuteState()
		if state.IsMute == nil || *state.IsMute != muted {
			t.Fatalf("wrong mute state: %+v", state)
		}
	}
	oscserver.Initialize(module, cfg)
	if oscserver.GetMuteState().IsMute != nil {
		t.Fatal("reinitialization must clear state")
	}
	cfg.Output.Osc.Ip = "[" // Malformed address fails without contacting a remote host.
	oscserver.Initialize(module, cfg)
	oscserver.HandleMessage(message.Message{Arg: message.MessageBody{Action: "mute"}})
	if oscserver.GetMuteState().IsMute != nil {
		t.Fatal("failed send must remain unknown")
	}
}
