package spec

import device "backend/dmxServer/devices"

func NewDimmer() *device.DMXDevice {
	return &device.DMXDevice{
		Model:      "dimmer",
		UseChannel: 1,
	}
}
