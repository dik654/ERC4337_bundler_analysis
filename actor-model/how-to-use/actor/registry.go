package actor

import (
	"sync"
)

const LocalLookupAddr = "local"

type Registry struct {
	mu     sync.RWMutex
	lookup map[string]Processer
	engine *Engine
}

// Registry 인스턴스 생성
// kind와 PID로 액터의 유형에 따른 관리
func newRegistry(e *Engine) *Registry {
	return &Registry{
		lookup: make(map[string]Processer, 1024),
		engine: e,
	}
}

// actor 종류와 PID를 이용하여 lookup table에서 Processor를 가져온 뒤
// Processor가 존재한다면 PID 리턴
func (r *Registry) GetPID(kind, id string) *PID {
	proc := r.getByID(kind + pidSeparator + id)
	if proc != nil {
		return proc.PID()
	}
	return nil
}

// lookup table에서 Processor 삭제
func (r *Registry) Remove(pid *PID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.lookup, pid.ID)
}

func (r *Registry) get(pid *PID) Processer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// PID에 해당하는 Processor를 가져와서 리턴
	if proc, ok := r.lookup[pid.ID]; ok {
		return proc
	}
	return nil
}

// lookup table에서 id에 해당하는 Processor 리턴
func (r *Registry) getByID(id string) Processer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lookup[id]
}

func (r *Registry) add(proc Processer) {
	r.mu.Lock()
	id := proc.PID().ID
	// Processor가 등록된 상태라면
	if _, ok := r.lookup[id]; ok {
		r.mu.Unlock()
		// 이미 존재하는 Actor를 추가하려했다는 Event 발생하고 종료
		r.engine.BroadcastEvent(ActorDuplicateIdEvent{PID: proc.PID()})
		return
	}
	// 아니라면 Processor 등록
	r.lookup[id] = proc
	r.mu.Unlock()
}
