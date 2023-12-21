package main

//pipeline data source
type DataSource interface {
	//return exit channel and data channel and
	//starts recieving data into data channel
	Start() (exit <-chan bool, data <-chan int)
	//stops recieving data and closes channels
	Stop()
}

//pipeline stage interface
type Stage interface {	
	//processes data from data channel and returns data into result channel. while processing checks exit channel
	Process(exit <-chan bool, data <-chan int) <-chan int
}

