package main

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

// cycle buffer drain interval 
const bufferDrainInterval time.Duration = 30 * time.Second

// buffer max capacity
const bufferSize int = 10


//issue:
//Стадия фильтрации чисел, не кратных 3 (не пропускать такие числа), исключая также и 0.
//Стадия буферизации данных в кольцевом буфере с интерфейсом, соответствующим тому,
//который был дан в качестве задания в 19 модуле. В этой стадии предусмотреть опустошение
//буфера (и соответственно, передачу этих данных, если они есть, дальше) с определённым
//интервалом во времени. Значения размера буфера и этого интервала времени
//сделать настраиваемыми (как мы делали: через константы или глобальные переменные).



func main() {
	log.SetOutput(io.Discard)//comment it to see logs

	//create pipeline
	pl := NewPipeline(
		NewConsoleDataSource(), // console data source (creates exit- and data- channels)
		NegativeFiltrationStage{}, // negative filter
		ModFiltrationStage{}, // mod 3 filter
		NewBufferizationStage(bufferSize, bufferDrainInterval), // cycle buffer
	)


	
	exit, res := pl.Run()

	//begin consumer
	wg := sync.WaitGroup{}	
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-exit:
				return
			case i := <-res:
				fmt.Println("Pipeline result: ", i)
			}
		}
	}()
	wg.Wait()
	//end consumer

	fmt.Println("Finish")

}
