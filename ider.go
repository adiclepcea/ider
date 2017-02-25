package ider

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

const max = 10000

//Ider is the main type
type Ider struct {
	MachineIndex uint
	curVal       uint
	chanIndex    chan int64
}

//NewIder will generate a new Ider with the specified settings
func NewIder(machineIndex uint) (*Ider, error) {
	if machineIndex > 9999 {
		return nil, fmt.Errorf("The maximum machineIndex for an Ider is 9999.")
	}

	ider := Ider{MachineIndex: machineIndex}

	ider.chanIndex = make(chan int64, 1)

	go ider.generateIndex()

	return &ider, nil
}

//GenerateID is used to get a new Id from ider
func (ider *Ider) GenerateID() int64 {

	c := <-ider.chanIndex

	//the channel might have one stall value in it and one waiting
	//so we get rid of those two if they are there
	for c/100000000 != time.Now().Unix() {
		c = <-ider.chanIndex
	}

	return c
}

//GenerateIDs is used to retrieve several ID's in seequence
//This should be use with care because of the limit of 9999 ids / second
func (ider *Ider) GenerateIDs(noOfIDs uint) []int64 {
	generatedIDs := make([]int64, noOfIDs)
	var i uint
	for i = 0; i < noOfIDs; i++ {
		generatedIDs[i] = ider.GenerateID()
	}
	return generatedIDs
}

func (ider *Ider) generateIndex() uint {
	var when int64
	for {
		now := time.Now().Unix()
		if now == when {
			if ider.curVal > max-2 {
				for now == when {
					now = time.Now().Unix()
					time.Sleep(1 * time.Millisecond)
				}
				ider.curVal = 1
			}
		} else {
			when = now
			ider.curVal = 1
		}

		ider.chanIndex <- when*int64(100000000) + int64(ider.curVal)*max + int64(ider.MachineIndex)
		ider.curVal++
	}
}
