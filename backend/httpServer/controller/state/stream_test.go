package state_test

import (
	"backend/config"
	dmxserver "backend/dmxServer"
	"backend/httpServer/controller/state"
	"backend/message"
	"backend/packageModule"
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/gin-gonic/gin"
)

func TestIdleStreamHeartbeat(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		response := httptest.NewRecorder()
		g, _ := gin.CreateTestContext(response)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		g.Request = httptest.NewRequestWithContext(ctx, http.MethodGet, "/stream", nil)
		done := make(chan struct{})
		go func() { defer close(done); state.Stream(g) }()
		synctest.Wait()
		initial := response.Body.String()
		if initial == "" {
			t.Fatal("missing initial snapshot")
		}
		time.Sleep(15 * time.Second)
		synctest.Wait()
		if response.Body.String() != initial+"\n" {
			t.Fatalf("expected one empty heartbeat, got %q", response.Body.String())
		}
		cancel()
		<-done
	})
}

func TestStreamSnapshotChangesAndCancellation(t *testing.T) {
	config.InitializeConfig()
	c := config.Get()
	c.Dmx.Groups = map[string]config.DMXGroup{"stage": {Devices: []config.DMXDevice{{Model: "dimmer", Channel: 1, MaxValue: []uint{255}}}}}
	module := packageModule.PackageModule{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), Wg: &sync.WaitGroup{}}
	if !dmxserver.Initialize(&module, &c) {
		t.Fatal("initialize failed")
	}
	defer dmxserver.CleanupDMXServer()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	done := make(chan struct{})
	engine.GET("/stream", func(g *gin.Context) { defer close(done); state.Stream(g) })
	server := httptest.NewServer(engine)
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/stream", nil)
	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.Header.Get("Content-Type") != "application/x-ndjson" || response.Header.Get("X-Accel-Buffering") != "no" {
		t.Fatal(response.Header)
	}
	frames := make(chan state.ControlState, 2)
	readErrors := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			if len(scanner.Bytes()) == 0 {
				continue
			}
			var frame state.ControlState
			if err := json.Unmarshal(scanner.Bytes(), &frame); err != nil {
				readErrors <- err
				return
			}
			frames <- frame
		}
		readErrors <- scanner.Err()
	}()
	next := func() state.ControlState {
		t.Helper()
		select {
		case frame := <-frames:
			return frame
		case err := <-readErrors:
			t.Fatalf("stream closed: %v", err)
		case <-time.After(time.Second):
			t.Fatal("missing state frame")
		}
		return state.ControlState{}
	}
	if frame := next(); frame.Fade["stage"].Level != 0 || frame.Mute.IsMute != nil {
		t.Fatal(frame)
	}
	select {
	case frame := <-frames:
		t.Fatalf("unchanged state was resent: %v", frame)
	case <-time.After(300 * time.Millisecond):
	}
	dmxserver.DMXServer.MessageHandler(message.Message{Arg: message.MessageBody{Action: "fade", Arg: map[string]string{"id": "stage", "isIn": "true", "duration": "0", "interval": "0"}}})
	dmxserver.Render()
	if frame := next(); frame.Fade["stage"].Level != 1 {
		t.Fatal(frame)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("disconnected stream did not stop")
	}
}
