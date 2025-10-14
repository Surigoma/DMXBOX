package dimmer

import device "backend/dmxServer/devices"

type Dimmer struct {
	device.DMXDevice
}

func (dim *Dimmer) Render() {

}
