package workerpool

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	URL string
}

type Result struct {
	URL          string
	Error        error
	StatusCode   int
	ResponseTime time.Duration
}

type Pool struct {
	stopped     bool
	workerCount int
	worker      *worker
	jobs        chan Job
	results     chan Result
	wg          *sync.WaitGroup
}

func (r Result) Info() string {
	if r.Error != nil {
		return fmt.Sprintf("[ERROR]: [%s]", r.URL)
	}
	return fmt.Sprintf("[SUCCESS]: [%s]", r.URL)
}

func (p *Pool) Stop() {
	p.stopped = true

	fmt.Printf("\nGRACEFUL SHUTDOWN:\n")

	close(p.jobs)

	p.wg.Wait()
}
