package rabbitmq

type Monitor struct {
	sent     int
	error    int
	received int
}

func (m *Monitor) IncSent() {
	m.sent++
}

func (m *Monitor) IncError() {
	m.error++
}

func (m *Monitor) IncReceived() {
	m.received++
}

func NewMonitor() *Monitor {
	return &Monitor{}
}
