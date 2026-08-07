package rabbtmq_go

type Monitor struct {
	sent  int
	error int
}

func (m *Monitor) IncSent() {
	m.sent++
}

func (m *Monitor) IncError() {
	m.error++
}
