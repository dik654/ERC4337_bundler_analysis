package main

import (
	"fmt"
	"sync"
)

// 함수 타입에 인터페이스 적용
type StoreProducerFunc func() Storer

type Storer interface {
	// getter, setter 인터페이스
	Push([]byte) (int, error)
	Get(int) ([]byte, error)
	Len() int
}

type MemoryStore struct {
	// 데이터를 안전하게 저장하기위해 mutex와 슬라이스
	mu   sync.RWMutex
	data [][]byte
}

// initialize
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		// bytes slices slice = [bytes slice][해당 bytes slice의 index에 해당하는 byte]
		data: make([][]byte, 0),
	}
}

// --methods--
func (s *MemoryStore) Push(b []byte) (int, error) {
	// RW lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// bytes slice 추가
	s.data = append(s.data, b)
	// index 리턴
	return len(s.data) - 1, nil
}

func (s *MemoryStore) Get(index int) ([]byte, error) {
	// R lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	// index 범위 체크
	if index < 0 {
		return nil, fmt.Errorf("offset cannot be smaller than 0")
	}
	if len(s.data)-1 < int(index) {
		return nil, fmt.Errorf("offset (%d) too high", index)
	}
	// 데이터 리턴
	return s.data[index], nil
}

func (s *MemoryStore) Len() int {
	return len(s.data)
}
