package fps

import "time"

type FPSController struct {
	FPS         float32
	recent      []float32
	recentIndex int
	Callback    func() bool
	Running     bool
	Finalize    func()
}

func NewFPS(fps float32, callback func() bool, finalize func()) *FPSController {
	if callback == nil {
		return nil
	}
	newFPS := FPSController{
		FPS:      fps,
		Callback: callback,
		Finalize: finalize,
	}
	newFPS.recent = make([]float32, 10)
	newFPS.recentIndex = 0
	return &newFPS
}

func (fps *FPSController) GetFPS() float32 {
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
	fps.Running = false
}

func (fps *FPSController) Run() {
	for i := range fps.recent {
		fps.recent[i] = 0
	}
	fps.recentIndex = 0
	privTime := time.Now()
	duration := time.Duration((1 / fps.FPS) * float32(time.Second))
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	fps.Running = true
	for fps.Running {
		<-ticker.C
		nextTime := time.Now()
		diff := nextTime.Sub(privTime)
		privTime = nextTime
		fps.recent[fps.recentIndex] = float32(1 / diff.Seconds())
		fps.recentIndex = (fps.recentIndex + 1) % len(fps.recent)
		if !fps.Callback() {
			break
		}
	}
	if fps.Finalize != nil {
		fps.Finalize()
	}
}
