package packet

type StatusRequest struct {
}

func (r *StatusRequest) ID() int {
	return 0x00
}

type StatusResponse struct {
	Status struct {
		Version struct {
			Name     string `json:"name"`
			Protocol int    `json:"protocol"`
		} `json:"version"`

		Players struct {
			Max    int `json:"max"`
			Online int `json:"online"`
		} `json:"players"`

		Description struct {
			Text string `json:"text"`
		} `json:"description"`
	}
}

func (r StatusResponse) ID() int {
	return 0x00
}

type StatusPing struct {
	Payload int64
}

func (p StatusPing) ID() int {
	return 0x01
}

type StatusPong struct {
	Payload int64
}

func (p StatusPong) ID() int {
	return 0x01
}
