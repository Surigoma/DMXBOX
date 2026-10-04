package dmxserver

import "time"

type FadeState struct {
	Level float64 `json:"level"`
	State string  `json:"state"`
	IsIn  bool    `json:"isIn"`
}

// GetFadeStates snapshots the software DMX output, including commands from TCP
// and other clients. It does not read back the physical fixtures.
func GetFadeStates() map[string]FadeState {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	result := make(map[string]FadeState, len(groups))
	now := time.Now()
	for id, group := range groups {
		state := FadeState{State: "idle"}
		var output, maximum float64
		for _, dev := range group.Devices {
			if dev == nil {
				continue
			}
			active, waiting, isIn := dev.FadeStatus(now)
			state.IsIn = isIn
			if active {
				if !waiting {
					state.State = "fading"
				} else if state.State != "fading" {
					state.State = "waiting"
				}
			}
			for i, max := range dev.MaxValue {
				if max == 0 {
					continue
				}
				value := (*dev.Output)[int(dev.Channel)-1+i]
				output += float64(min(value, max))
				maximum += float64(max)
			}
		}
		if maximum > 0 {
			state.Level = output / maximum
		}
		result[id] = state
	}
	return result
}
