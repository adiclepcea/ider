package ider

import (
	"fmt"
	"sync"
	"time"
)

var chanIndex chan uint
var wg sync.WaitGroup

const max = 10000

//Ider is the main type
type Ider struct {
	MachineIndex uint
	curVal       uint
	chanIndex    chan uint
}

//NewIder will generate a new Ider with the specified settings
func NewIder(machineIndex uint) (*Ider, error) {
	if machineIndex > 9999 {
		return nil, fmt.Errorf("The maximum machineIndex for an Ider is 9999.")
	}

	ider := Ider{MachineIndex: machineIndex}

	ider.chanIndex = make(chan uint, 1)

	go ider.generateIndex()

	return &ider, nil
}

//GenerateID is used to get a new Id from ider
func (ider *Ider) GenerateID() int64 {
	c := <-ider.chanIndex
	return time.Now().Unix()*int64(100000000) + int64(c)*max + int64(ider.MachineIndex)
}

func (ider *Ider) generateIndex() uint {
	var when int64
	for {
		now := time.Now().Unix()
		if now == when {
			if ider.curVal > max-2 {
				time.Sleep(1 * time.Millisecond)
			}
		} else {
			when = now
			ider.curVal = 0
		}
		ider.curVal++
		ider.chanIndex <- ider.curVal
	}
}
