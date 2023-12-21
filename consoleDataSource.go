package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type consoleDataSource struct {
	exit    chan bool   //exit channel
	console chan int    //data channel
	stopped bool        //flag if stopped then true
	started bool        //flag if started then true
	m       *sync.Mutex //lock for flags
}

// creates new Console datasource
func NewConsoleDataSource() DataSource {
	return &consoleDataSource{m: &sync.Mutex{}}
}

// returns data and exit channels and
// starts recieving data from console
func (cds *consoleDataSource) Start() (<-chan bool, <-chan int) {

	//if started then exiting
	cds.m.Lock()
	if cds.started && !cds.stopped {
		defer cds.m.Unlock()
		return cds.exit, cds.console
	}
	cds.m.Unlock()

	//welcome text
	fmt.Println("Hi there! Console Data Source started to scan.")
	fmt.Println("To exit type 'exit'.")
	fmt.Println("Type integer and then press 'enter' to process it.")

	//initiate channels
	cds.exit = make(chan bool)
	cds.console = make(chan int)

	//mark 'started'
	cds.m.Lock()
	cds.started = true
	cds.stopped = false
	cds.m.Unlock()

	go func() {
		defer cds.Stop()

		scanner := bufio.NewScanner(os.Stdin)

		for {
			//check if there is external stop
			select {
			case <-cds.exit:
				fmt.Println("External 'exit' command recieved. Breaking...")
				return
			default:
			}

			//reading console and check data
			scanner.Scan()
			consoleData := scanner.Text()
			if strings.EqualFold(consoleData, "exit") {
				fmt.Println("'Exit' command recieved. Breaking...")
				return
			}
			i, err := strconv.Atoi(consoleData)
			if err != nil {
				fmt.Println("Error. Enter integer number.")
				continue
			}

			//output correct data
			cds.console <- i
		}
	}()

	return cds.exit, cds.console
}

// stops recieving data and closes channels
func (cds *consoleDataSource) Stop() {
	//if stopped - exit
	cds.m.Lock()
	defer cds.m.Unlock()
	if !cds.started || cds.stopped {
		return
	}
	cds.m.Unlock()

	close(cds.exit)
	close(cds.console)

	//mark 'stopped'
	cds.m.Lock()
	cds.stopped = true
	cds.started = false

}
