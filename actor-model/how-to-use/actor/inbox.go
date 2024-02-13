package actor

import (
	"runtime"
	"sync/atomic"

	"github.com/anthdm/hollywood/ringbuffer"
)

const (
	defaultThroughput = 300
	messageBatchSize  = 1024 * 4
)

const (
	idle int32 = iota
	running
	stopped
)

type Scheduler interface {
	Schedule(fn func())
	Throughput() int
}

// int 타입을 goscheduler라는 타입으로 재정의
type goscheduler int

func (goscheduler) Schedule(fn func()) {
	go fn()
}

func (sched goscheduler) Throughput() int {
	// 처리율을 int값으로 다시 리턴
	return int(sched)
}

// 재정의한 타입으로 처리율 int값을 변환
func NewScheduler(throughput int) Scheduler {
	return goscheduler(throughput)
}

type Inboxer interface {
	Send(Envelope)
	Start(Processer)
	Stop() error
}

// actor가 받는 메세지를 저장하는 버퍼
type Inbox struct {
	// 데이터 저장하는 링 버퍼
	rb *ringbuffer.RingBuffer[Envelope]
	// Processor(actor)
	proc Processer
	// 처리율과 고루틴 생성을 하는 scheduler 구조체 저장
	scheduler Scheduler
	// 현재 Processor의 상태를 저장
	procStatus int32
}

// initialize
func NewInbox(size int) *Inbox {
	return &Inbox{
		rb:        ringbuffer.New[Envelope](int64(size)),
		scheduler: NewScheduler(defaultThroughput),
	}
}

func (in *Inbox) Send(msg Envelope) {
	in.rb.Push(msg)
	in.schedule()
}

func (in *Inbox) schedule() {
	// atomic은 동시성을 안전하게 처리
	// Inbox의 procStatus가 현재 idle 상태일 때만 running 상태로 변경
	if atomic.CompareAndSwapInt32(&in.procStatus, idle, running) {
		// process 메서드 참조를 전달하여 go routine 실행
		in.scheduler.Schedule(in.process)
	}
}

func (in *Inbox) process() {
	in.run()
	// 실행이 완료되면 running 상태를 idle 상태로 변결
	atomic.StoreInt32(&in.procStatus, idle)
}

func (in *Inbox) run() {
	// i, t를 0으로 초기화하고
	// t에 처리량값 저장
	i, t := 0, in.scheduler.Throughput()
	// 현재 processor의 상태가 stopped가 아니라면
	for atomic.LoadInt32(&in.procStatus) != stopped {
		// 작업 수가 처리량을 넘어선다면
		if i > t {
			i = 0
			// 현재 고루틴이 CPU를 양보하여 다른 고루틴에게 실행 기회를 줌
			runtime.Gosched()
		}
		i++

		// 링 버퍼에서 메세지를 꺼내서 ReceiveFunc 실행
		if msgs, ok := in.rb.PopN(messageBatchSize); ok && len(msgs) > 0 {
			in.proc.Invoke(msgs)
		} else {
			return
		}
	}
}

// 수신함에 actor 할당
func (in *Inbox) Start(proc Processer) {
	in.proc = proc
}

// 수신함의 상태를 stopped로 변경
func (in *Inbox) Stop() error {
	atomic.StoreInt32(&in.procStatus, stopped)
	return nil
}
