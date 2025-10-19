package spec

import device "backend/dmxServer/devices"

func NewWCLight() *device.DMXDevice {
	return &device.DMXDevice{
		Model:      "wclight",
		UseChannel: 3,
	}
}
