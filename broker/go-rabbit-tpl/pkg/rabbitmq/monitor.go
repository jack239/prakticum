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

func (m *Monitor) Stats() map[string]int {
	return map[string]int{
		"sent":     m.sent,
		"error":    m.error,
		"received": m.received,
	}
}

func NewMonitor() *Monitor {
	return &Monitor{}
}
