package seal

import "sync"

type Manager struct {
	sealed bool
	key    []byte
	mu     sync.RWMutex
}

func NewManger() *Manager {
	return &Manager{
		sealed: true,
	}
}

func (m *Manager) IsSealed() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sealed
}

func (m *Manager) Seal() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.key = nil
	m.sealed = true
}

func (m *Manager) Unseal(key []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.key = key
	m.sealed = false
}

func (m *Manager) GetKey() []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.key
}
