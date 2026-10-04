package state

import (
	dmxserver "backend/dmxServer"
	oscserver "backend/oscServer"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ControlState struct {
	Fade map[string]dmxserver.FadeState `json:"fade"`
	Mute oscserver.MuteState            `json:"mute"`
}

// Stream sends an initial snapshot and changed snapshots as newline-delimited JSON.
// Empty lines keep idle connections alive without generating HTTP access requests.
func Stream(g *gin.Context) {
	g.Header("Content-Type", "application/x-ndjson")
	g.Header("Cache-Control", "no-cache, no-transform")
	g.Header("X-Accel-Buffering", "no")
	writer := http.NewResponseController(g.Writer)
	tick := time.NewTicker(250 * time.Millisecond)
	heartbeat := time.NewTicker(15 * time.Second)
	defer tick.Stop()
	defer heartbeat.Stop()
	var previous []byte
	write := func(body []byte) bool {
		// Bound slow clients so shutdown and disconnected readers cannot leave writers stuck.
		_ = writer.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if _, err := g.Writer.Write(append(body, '\n')); err != nil {
			return false
		}
		return writer.Flush() == nil
	}
	snapshot := func() bool {
		body, err := json.Marshal(ControlState{Fade: dmxserver.GetFadeStates(), Mute: oscserver.GetMuteState()})
		if err != nil {
			return false
		}
		if bytes.Equal(body, previous) {
			return true
		}
		if !write(body) {
			return false
		}
		previous = body
		return true
	}
	if !snapshot() {
		return
	}
	for {
		select {
		case <-g.Request.Context().Done():
			return
		case <-tick.C:
			if !snapshot() {
				return
			}
		case <-heartbeat.C:
			if !write(nil) {
				return
			}
		}
	}
}
