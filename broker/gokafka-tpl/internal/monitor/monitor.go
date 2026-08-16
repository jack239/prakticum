package monitor

import "strconv"

type Monitor struct {
	errors    int
	published int
}

func (m *Monitor) IncError() {
	m.errors++
}

func (m *Monitor) IncPublished() {
	m.published++
}

func (m *Monitor) Stats() map[string]string {
	return map[string]string{
		"errors":    strconv.Itoa(m.errors),
		"published": strconv.Itoa(m.published),
	}
}

func NewMonitor() *Monitor {
	return &Monitor{}
}
