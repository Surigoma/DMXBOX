package device

import (
	"math"
	"sync"
	"time"
)

type DMXDevice struct {
	Model       string
	Channel     uint8
	UseChannel  uint8
	Output      *[]byte
	Before      []byte
	Target      []byte
	MaxValue    []byte
	effectStart time.Time
	effectEnd   time.Time
	Duration    *float32
	once        bool
	ModFade     func(isIn bool, duration float32, interval float32)
	ModUpdate   func() bool
}

func (dev *DMXDevice) Initialize(channel uint8, maxValue []byte, target *[]byte, duration *float32) bool {
	dev.Channel = channel
	dev.Output = target
	dev.Duration = duration
	dev.Before = make([]byte, dev.UseChannel)
	dev.MaxValue = make([]byte, dev.UseChannel)
	dev.Target = make([]byte, dev.UseChannel)
	dev.effectStart = time.Now()
	dev.effectEnd = time.Now()
	if len(dev.MaxValue) > len(maxValue) {
		return false
	}
	for i := range dev.MaxValue {
		dev.MaxValue[i] = maxValue[i]
	}
	for i := range *dev.Output {
		(*dev.Output)[i] = 0
	}
	for i := range dev.Target {
		dev.Target[i] = 0
	}
	return true
}

func (dev *DMXDevice) Fade(isIn bool, duration float32, interval float32) {
	if dev.ModFade != nil {
		dev.ModFade(isIn, duration, interval)
		return
	}
	dur := *dev.Duration
	inter := float32(0)
	if duration >= 0 {
		dur = duration
	}
	if interval > 0 {
		inter = interval
	}
	dev.effectStart = time.Now().Add(time.Duration(inter * float32(time.Second)))
	dev.effectEnd = dev.effectStart.Add(time.Duration(dur * float32(time.Second)))
	if dur == 0 {
		dev.once = true
	}

	for i := range dev.Before {
		dev.Before[i] = (*dev.Output)[i+int(dev.Channel)-1]
	}
	for i := range dev.Target {
		if isIn {
			dev.Target[i] = dev.MaxValue[i]
		} else {
			dev.Target[i] = 0
		}
	}
}

func (dev *DMXDevice) Update(wg *sync.WaitGroup) bool {
	defer wg.Done()
	if dev.ModUpdate != nil {
		return dev.ModUpdate()
	}
	now := time.Now()
	nowD := now.Sub(dev.effectStart)
	endD := dev.effectEnd.Sub(dev.effectStart)
	percentRaw := nowD.Seconds() / endD.Seconds()
	percent := math.Max(0.0, math.Min(1.0, percentRaw))
	if nowD < 0 {
		return false
	}
	if dev.once {
		dev.once = false
		for i := range dev.Target {
			(*dev.Output)[i+int(dev.Channel)-1] = dev.Target[i]
		}
		return true
	}
	if percentRaw > 1.1 {
		return false
	}
	for i := range dev.Target {
		v := (float64(dev.Target[i])-float64(dev.Before[i]))*float64(percent) + float64(dev.Before[i])
		(*dev.Output)[i+int(dev.Channel)-1] = byte(math.Max(0, math.Min(255, math.Round(v))))
	}
	return true
}
