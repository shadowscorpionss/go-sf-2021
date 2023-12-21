package main

type Pipeline struct {
	ds     DataSource
	stages []Stage
}
//creates pipeline of stages
func NewPipeline(ds DataSource, stages ...Stage) *Pipeline {
	return &Pipeline{ds: ds, stages: stages}
}

//starts pipeline
func (p *Pipeline) Run() (<-chan bool, <-chan int) {
	exit, res := p.ds.Start()
	for _, stg := range p.stages {
		res = stg.Process(exit, res)
	}
	return exit, res
}

//breaks pipeline (if needed)
func (p *Pipeline) Break() {
	p.ds.Stop()
}
