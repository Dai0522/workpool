package workpool

import (
	"encoding/binary"
)

// SampleTask .
type SampleTask struct {
	ID uint64
}

// Run .
func (t *SampleTask) Run() *[]byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, t.ID)
	return &result
}
