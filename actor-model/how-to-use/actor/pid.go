package actor

import (
	"github.com/zeebo/xxh3"
)

const pidSeparator = "/"

// actor.proto
//
//	message PID {
//		string address = 1;
//		string ID = 2;
//	}
func NewPID(address, id string) *PID {
	p := &PID{
		Address: address,
		// id는 kind + pidSeparator + ID로 구성
		ID: id,
	}
	return p
}

// PID getter
func (pid *PID) String() string {
	return pid.Address + pidSeparator + pid.ID
}

// 동일한 PID인지 체크
func (pid *PID) Equals(other *PID) bool {
	return pid.Address == other.Address && pid.ID == other.ID
}

// 부모와 동일한 address로 자식 PID 생성
func (pid *PID) Child(id string) *PID {
	childID := pid.ID + pidSeparator + id
	return NewPID(pid.Address, childID)
}

// 해당 PID의 해시 리턴
func (pid *PID) LookupKey() uint64 {
	key := []byte(pid.Address)
	key = append(key, pid.ID...)
	return xxh3.Hash(key)
}
