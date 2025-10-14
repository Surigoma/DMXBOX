package device

type DMXDevice struct {
	Channel    uint8
	useChannel uint8
	target     *[]byte
	Render     func(target *[]byte) bool
}

func (dev *DMXDevice) Initialize(channel uint8) bool {
	dev.Channel = channel
	return false
}

func (dev *DMXDevice) Update() bool {
	return dev.Render(dev.target)
}
