package config

import (
	"fmt"
	"math"
	"strings"
)

// Validate rejects settings that cannot be used safely by the selected modules.
func (c Config) Validate() error {
	finite := func(v float32) bool { return !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) }
	if !finite(c.Dmx.Fps) || c.Dmx.Fps <= 0 || c.Dmx.Fps > 1000 {
		return fmt.Errorf("dmx.fps must be between 0 and 1000 (exclusive of 0)")
	}
	if !finite(c.Dmx.FadeInterval) || c.Dmx.FadeInterval < 0 || !finite(c.Dmx.Delay) || c.Dmx.Delay < 0 {
		return fmt.Errorf("fadeInterval and delay must be finite, non-negative values")
	}
	for _, m := range c.Input.Modules {
		switch m {
		case "http":
			if c.Input.Http.Port == 0 {
				return fmt.Errorf("HTTP port must be non-zero")
			}
		case "tcp":
			if c.Input.Tcp.Port == 0 {
				return fmt.Errorf("TCP port must be non-zero")
			}
		default:
			return fmt.Errorf("unsupported input module: %s", m)
		}
	}
	if len(c.Output.Target) == 0 {
		return fmt.Errorf("select at least one output")
	}
	for _, m := range c.Output.Target {
		switch m {
		case "console", "ftdi":
		case "artnet":
			if c.Output.Artnet.Universe > 15 || c.Output.Artnet.SubUniverse > 15 || c.Output.Artnet.Net > 127 {
				return fmt.Errorf("invalid Art-Net universe, subnet or net")
			}
			if p := c.Output.Artnet.Port; p != nil && (*p == 0 || *p > 65535) {
				return fmt.Errorf("invalid Art-Net port")
			}
		case "osc":
			if c.Output.Osc.Type != "int" && c.Output.Osc.Type != "float" {
				return fmt.Errorf("OSC type must be int or float")
			}
			if c.Output.Osc.Port == 0 || !strings.HasPrefix(c.Output.Osc.Format, "/") {
				return fmt.Errorf("invalid OSC port or address format")
			}
		default:
			return fmt.Errorf("unsupported output module: %s", m)
		}
	}
	for id, group := range c.Dmx.Groups {
		if strings.TrimSpace(id) == "" || strings.ContainsAny(id, ".[]") || id == "__proto__" || id == "constructor" || id == "prototype" {
			return fmt.Errorf("invalid group ID: %q", id)
		}
		for i, d := range group.Devices {
			width := 0
			switch d.Model {
			case "dimmer":
				width = 1
			case "wclight":
				width = 3
			default:
				return fmt.Errorf("group %s device %d: unsupported model %s", id, i, d.Model)
			}
			if d.Channel == 0 || int(d.Channel)+width-1 > 512 {
				return fmt.Errorf("group %s device %d: channels must fit within 1..512", id, i)
			}
			if len(d.MaxValue) < width {
				return fmt.Errorf("group %s device %d: missing max values", id, i)
			}
			for _, v := range d.MaxValue {
				if v > 255 {
					return fmt.Errorf("group %s device %d: max must be within 0..255", id, i)
				}
			}
		}
	}
	return nil
}
