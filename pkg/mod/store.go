package mod

import (
	"fmt"
	"sync"
	"time"
)

// Task 是共享资源，v1 和 v2 操作的是同一个 Store 实例

type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"` // pending / running / done
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store struct {
	mu    sync.RWMutex
	tasks map[string]*Task
	seq   int
}

func NewStore() *Store {
	return &Store{tasks: make(map[string]*Task)}
}

func (s *Store) List() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

func (s *Store) Get(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	return t, nil
}

func (s *Store) Create(name string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	now := time.Now()
	t := &Task{
		ID:        fmt.Sprintf("task-%03d", s.seq),
		Name:      name,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.tasks[t.ID] = t
	return t
}

func (s *Store) Update(id, name, status string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	if name != "" {
		t.Name = name
	}
	if status != "" {
		t.Status = status
	}
	t.UpdatedAt = time.Now()
	return t, nil
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found")
	}
	delete(s.tasks, id)
	return nil
}
