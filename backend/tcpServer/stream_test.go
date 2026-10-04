package tcpserver

import (
	"backend/config"
	"backend/packageModule"
	"bufio"
	"io"
	"log/slog"
	"net"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCommandStreamFragments(t *testing.T) {
	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()
	commands := make(chan string, 4)
	done := make(chan error, 1)
	go func() {
		s := bufio.NewScanner(reader)
		s.Buffer(make([]byte, 512), 4098)
		s.Split(splitCommand)
		for s.Scan() {
			if text := s.Text(); text != "" {
				commands <- text
			}
		}
		done <- s.Err()
	}()
	writer.Write([]byte("fadeIn st"))
	select {
	case command := <-commands:
		t.Fatalf("partial command dispatched: %q", command)
	default:
	}
	writer.Write([]byte("g duration:0,interval:0\r"))
	select {
	case command := <-commands:
		if command != "fadeIn stg duration:0,interval:0" {
			t.Fatal(command)
		}
	case <-time.After(time.Second):
		t.Fatal("CR-terminated command was not dispatched")
	}
	writer.Write([]byte("\nfadeOut stg\nmute false\r\n"))
	// A command can span multiple 512-byte reads.
	long := "fadeIn " + strings.Repeat("a", 700)
	writer.Write([]byte(long + "\n"))
	writer.Close()
	var got []string
	for range 3 {
		got = append(got, <-commands)
	}
	if want := []string{"fadeOut stg", "mute false", long}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestCommandStreamLengthLimit(t *testing.T) {
	for _, length := range []int{4095, 4096} {
		s := bufio.NewScanner(strings.NewReader(strings.Repeat("a", length) + "\r\n"))
		s.Buffer(make([]byte, 512), 4098)
		s.Split(splitCommand)
		if got := s.Scan(); got != (length == 4095) {
			t.Fatalf("length %d: accepted = %v", length, got)
		}
		if length == 4096 && s.Err() == nil {
			t.Fatal("missing length error")
		}
	}
}

func TestStopClosesPendingConnections(t *testing.T) {
	config.InitializeConfig()
	c := config.Get()
	c.Input.Tcp.Port = 0 // Let the OS allocate an isolated test port.
	var moduleWg sync.WaitGroup
	module := packageModule.PackageModule{ModuleName: "tcp", Wg: &moduleWg, Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	if !Initialize(&module, &c) {
		t.Fatal("initialize failed")
	}
	moduleWg.Add(1)
	StartTCP()
	listenerMutex.Lock()
	address := listener.Addr().String()
	listenerMutex.Unlock()
	conn, err := net.Dial("tcp", address)
	if err != nil {
		StopTCP()
		t.Fatal(err)
	}
	defer conn.Close()
	conn.Write([]byte("fadeIn"))
	StopTCP()
	moduleWg.Wait()
	conn.SetReadDeadline(time.Now().Add(time.Second))
	if _, err := conn.Read(make([]byte, 1)); err == nil {
		t.Fatal("connection remained open")
	} else if timeout, ok := err.(net.Error); ok && timeout.Timeout() {
		t.Fatal("connection was not closed on stop")
	}
}
