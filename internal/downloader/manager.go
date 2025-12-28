package downloader

import (
	"context"
	"sync"
)

type Manager struct {
	queue    chan Job
	updates  chan Progress
	contexts map[int]context.CancelFunc
	mu       sync.Mutex
}

func NewManager(workers int) *Manager {
	m := &Manager{
		queue:    make(chan Job, 100),
		updates:  make(chan Progress, 100),
		contexts: make(map[int]context.CancelFunc),
	}

	for i := 0; i < workers; i++ {
		go func() {
			for job := range m.queue {
				ctx, cancel := context.WithCancel(context.Background())
				m.mu.Lock()
				m.contexts[job.ID] = cancel
				m.mu.Unlock()
				job.Run(ctx, m.updates)
				m.mu.Lock()
				delete(m.contexts, job.ID)
				m.mu.Unlock()
			}
		}()
	}

	return m
}

func (m *Manager) Submit(job Job) {
	m.queue <- job
}

func (m *Manager) Updates() <-chan Progress {
	return m.updates
}

func (m *Manager) Cancel(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cancel, ok := m.contexts[id]; ok {
		cancel()
		delete(m.contexts, id)
	}
}
