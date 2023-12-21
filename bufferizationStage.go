package main

import (
	"log"
	"time"
)

type BufferizationStage struct {
	cap           int           //capacity of buffer
	drainInterval time.Duration //buffer drain interval
	buffer        *CycleIntBuffer
}

func NewBufferizationStage(capacity int, drainInterval time.Duration) Stage {
	return &BufferizationStage{
		cap:           capacity,
		buffer:        NewCycleIntBuffer(capacity),
		drainInterval: drainInterval,
	}

}

func (bs *BufferizationStage) Process(exit <-chan bool, data <-chan int) <-chan int {
	res := make(chan int)

	buffer := bs.buffer //for shorter calls :)
	buffer.Clear()

	go func() {

		for {
			select {
			case <-exit:
				log.Println("bs: exit. breaking...")
				return
			case i, isChannelOpen := <-data:
				if !isChannelOpen {
					log.Println("bs: data channel is closed. breaking...")
					return
				}
				//buffering
				buffer.Push(i)
				log.Printf("bs: +%d (%d)\n", i, buffer.Count())

			}
		}

	}()

	go func() {
		defer close(res)
		for {
			//sending buffered
			select {

			case <-exit:
				log.Println("bs: exit. breaking...")

				return
			case <-time.After(bs.drainInterval):
				da := buffer.Get()
				log.Print("bs: checking buffer data... ")

				if da != nil {
					for _, d := range da {
						select {
						case <-exit:
							log.Println("bs: exit. breaking...")

							return
						case res <- d:
							log.Printf("bs: -> %d \n", d)
						}
					}
				} else {
					log.Println("nil")
				}

			}

		}
	}()

	return res

}
