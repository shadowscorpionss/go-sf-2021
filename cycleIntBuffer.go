package main

import (
	"sync"
)

// buffer of integers
// when capacity is discovered moves all items to begin
// then pushes item to end. And so on in cycle
type CycleIntBuffer struct {
	arr []int 		// int array for cycle buffer 
	pos  int        // current position
	cap int    		// capacity of buffer
	m    sync.Mutex // lock object for thread-safe push and get

}

// creates cycle buffer of integers
func NewCycleIntBuffer(capacity int) *CycleIntBuffer {
	return &CycleIntBuffer{make([]int, capacity), -1, capacity, sync.Mutex{}}
}

// adds item to the end of buffer. 
// when capacity is discovered new items replaces old by 
// moving to begin of buffer
func (cib *CycleIntBuffer) Push(item int) {
	cib.m.Lock()
	defer cib.m.Unlock()
	//check if capacity is discovered
	if cib.pos == cib.cap-1 {
		//is discovered: moving items to begin
		for i := 1; i <= cib.cap-1; i++ {
			cib.arr[i-1] = cib.arr[i]
		}
	} else {
		//is not discoverd. incrementing position
		cib.pos++	
	}
	//set last item (item at positionß)
	cib.arr[cib.pos] = item
}

// returns all items in buffer and cleans buffer
func (cib *CycleIntBuffer) Get() []int {
	cib.m.Lock()
	if cib.pos < 0 {
		defer cib.m.Unlock()
		return nil
	}

	res := cib.arr[:cib.pos+1]
	cib.m.Unlock()

	cib.Clear()
	return res
}

//clears items
func (cib *CycleIntBuffer) Clear(){
	cib.m.Lock()
	defer cib.m.Unlock()
	cib.pos=-1
}

//returns count of items in buffer
func (cib *CycleIntBuffer) Count()int{
	cib.m.Lock()
	defer cib.m.Unlock()
	return cib.pos+1
}
