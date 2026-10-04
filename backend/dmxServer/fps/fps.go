package fps

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type FPSController struct {
	FPS         float32
	recent      []float32
	recentIndex int
	Callback    func() bool
	Running     atomic.Bool
	lock        sync.Mutex
	stop        chan struct{}
	stopOnce    sync.Once
	Finalize    func()
}

func NewFPS(fps float32, callback func() bool, finalize func()) *FPSController {
	if callback == nil || fps <= 0 || math.IsNaN(float64(fps)) || math.IsInf(float64(fps), 0) || time.Duration(float64(time.Second)/float64(fps)) <= 0 {
		return nil
	}
	newFPS := FPSController{
		FPS:      fps,
		Callback: callback,
		Finalize: finalize,
		stop:     make(chan struct{}),
	}
	newFPS.recent = make([]float32, 10)
	newFPS.recentIndex = 0
	return &newFPS
}

func (fps *FPSController) GetFPS() float32 {
	fps.lock.Lock()
	defer fps.lock.Unlock()
	var sum float32 = 0
	length := 0
	for _, v := range fps.recent {
		if v == 0 {
			continue
		}
		sum += v
		length++
	}
	if length <= 0 {
		return -1
	}
	return sum / float32(length)
}

func (fps *FPSController) Stop() {
	fps.stopOnce.Do(func() { close(fps.stop) })
	fps.Running.Store(false)
}

func (fps *FPSController) Run() {
	fps.lock.Lock()
	for i := range fps.recent {
		fps.recent[i] = 0
	}
	fps.recentIndex = 0
	fps.lock.Unlock()
	privTime := time.Now()
	duration := time.Duration((1 / fps.FPS) * float32(time.Second))
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	fps.Running.Store(true)
	defer fps.Running.Store(false)
loop:
	for {
		select {
		case <-fps.stop:
			break loop
		case <-ticker.C:
		}
		nextTime := time.Now()
		diff := nextTime.Sub(privTime)
		privTime = nextTime
		fps.lock.Lock()
		fps.recent[fps.recentIndex] = float32(1 / diff.Seconds())
		fps.recentIndex = (fps.recentIndex + 1) % len(fps.recent)
		fps.lock.Unlock()
		if !fps.Callback() {
			break
		}
	}
	fps.Running.Store(false)
	if fps.Finalize != nil {
		fps.Finalize()
	}
}
