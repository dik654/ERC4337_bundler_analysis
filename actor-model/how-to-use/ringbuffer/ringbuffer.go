package ringbuffer

import (
	"sync"
	"sync/atomic"
)

type buffer[T any] struct {
	items           []T
	head, tail, mod int64
}

type RingBuffer[T any] struct {
	len     int64
	content *buffer[T]
	mu      sync.Mutex
}

func New[T any](size int64) *RingBuffer[T] {
	return &RingBuffer[T]{
		content: &buffer[T]{
			items: make([]T, size),
			head:  0,
			tail:  0,
			mod:   size,
		},
		len: 0,
	}
}

func (rb *RingBuffer[T]) Push(item T) {
	rb.mu.Lock()
	rb.content.tail = (rb.content.tail + 1) % rb.content.mod
	if rb.content.tail == rb.content.head {
		size := rb.content.mod * 2
		newBuff := make([]T, size)
		for i := int64(0); i < rb.content.mod; i++ {
			idx := (rb.content.tail + i) % rb.content.mod
			newBuff[i] = rb.content.items[idx]
		}
		content := &buffer[T]{
			items: newBuff,
			head:  0,
			tail:  rb.content.mod,
			mod:   size,
		}
		rb.content = content
	}
	atomic.AddInt64(&rb.len, 1)
	rb.content.items[rb.content.tail] = item
	rb.mu.Unlock()
}

func (rb *RingBuffer[T]) Len() int64 {
	return atomic.LoadInt64(&rb.len)
}

func (rb *RingBuffer[T]) Pop() (T, bool) {
	if rb.Len() == 0 {
		var t T
		return t, false
	}
	rb.mu.Lock()
	rb.content.head = (rb.content.head + 1) % rb.content.mod
	item := rb.content.items[rb.content.head]
	var t T
	rb.content.items[rb.content.head] = t
	atomic.AddInt64(&rb.len, -1)
	rb.mu.Unlock()
	return item, true
}

func (rb *RingBuffer[T]) PopN(n int64) ([]T, bool) {
	// 링 버퍼의 길이가 0이라면 종료
	if rb.Len() == 0 {
		return nil, false
	}
	// mutex RW lock 걸기
	rb.mu.Lock()
	content := rb.content

	// 꺼내려는 개수가 링 버퍼의 길이보다 크면 총 길이로 설정
	if n >= rb.len {
		n = rb.len
	}
	atomic.AddInt64(&rb.len, -n)

	// 꺼내려는 개수만큼의 크기의 임시 슬라이스 생성
	items := make([]T, n)
	for i := int64(0); i < n; i++ {
		// 임시 슬라이스에 버퍼의 데이터를 담고
		pos := (content.head + 1 + i) % content.mod
		items[i] = content.items[pos]
		// 버퍼 비우기
		var t T
		content.items[pos] = t
	}
	// 꺼낸 마지막 값의 다음 값이 링 버퍼의 시작 위치
	content.head = (content.head + n) % content.mod

	// lock 풀기
	rb.mu.Unlock()
	// 임시 슬라이스 리턴
	return items, true
}
