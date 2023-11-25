package workerpool

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type worker struct {
	client *http.Client
}

func newWorker(timeout time.Duration) *worker {
	return &worker{
		&http.Client{
			Timeout: timeout,
		},
	}
}

func New(workerCount int, timeout time.Duration, result chan Result) *Pool {
	return &Pool{
		worker:      newWorker(timeout),
		results:     result,
		cntJob:      new(sync.WaitGroup),
		cntResult:   new(sync.WaitGroup),
		workerCount: workerCount,
		jobs:        make(chan Job),
		stopped:     false,
	}
}

func (w worker) process(j Job) Result {
	result := Result{URL: j.URL}

	now := time.Now()

	resp, err := w.client.Get(j.URL)
	if err != nil {
		result.Error = err
		return result
	}
	result.StatusCode = resp.StatusCode
	result.ResponseTime = time.Since(now)
	return result
}

func RunWorker(worker *Pool) { // ПЕРВЫЙ
	for i := 0; i < WORKERS_COUNT; i++ {
		go worker.GetWorker(i)
	}
}

func (worker *Pool) GetWorker(ID int) { // ВТОРОЙ
	for job := range worker.jobs {
		time.Sleep(time.Second)
		worker.cntResult.Add(1)
		worker.results <- worker.worker.process(job)
		worker.cntJob.Done()
	}
	log.Printf("worker ID %d finished", ID)
}

func (worker *Pool) PushURL(urls []string) { // ТРЕТИЙ
	for {
		for _, url := range urls {
			if worker.stopped {
				return
			}
			worker.jobs <- Job{URL: url}
			worker.cntJob.Add(1)
		}
		time.Sleep(INTERVAL)
		fmt.Println("---------------")
	}
}

func (worker *Pool) GetResult() { // ЧЕТВЕРТЫЙ
	for result := range worker.results {
		fmt.Println(result.Info())
		worker.cntResult.Done()
	}
}
