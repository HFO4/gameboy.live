package driver

type ChannelIO struct {
	Data   byte
	Target *ChannelIO
	Open   bool
	Master bool

	SendDelay int

	Receive chan byte
}

func (io *ChannelIO) SetTarget(p *ChannelIO) {
	io.Target = p
}

func (io *ChannelIO) SetChannelStatus(master bool, status bool) {
	io.Open = status
	io.Master = master
}

func (io *ChannelIO) SendByte(data byte) bool {
	if io.SendDelay != 0 {
		//log.Fatal("delay,",io.SendDelay)
	}
	io.Data = data
	if io.Master {
		io.SendDelay = 4000
	}

	return false
}

func (io *ChannelIO) FetchByte(cycles int) (byte, bool) {
	if io.SendDelay == 0 {
		select {
		case data := <-io.Receive:
			if io.Target != nil {
				io.Target.Receive <- io.Data
				return data, true
			} else {
				return 0xff, false
			}

		default:
			return 0xff, false
		}

	}

	io.SendDelay -= cycles

	if io.SendDelay <= 0 {

		if io.Target == nil {

			return 0xff, true
		}

		io.Target.Receive <- io.Data
		received := <-io.Receive
		io.SendDelay = 0
		return received, true
	}

	return 0xff, false
}
