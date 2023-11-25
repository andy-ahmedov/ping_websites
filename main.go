package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/andy-ahmedov/ping_websites/workerpool"
	_ "github.com/andy-ahmedov/ping_websites/workerpool"
)

var urls = []string{
	"https://ya.ru/",
	"https://google.ru/",
	"https://google.com/",
	"https://www.canva.com",
	"https://kinoland.biz/",
}

func main() {
	result := make(chan workerpool.Result)

	workerPool := workerpool.New(workerpool.WORKERS_COUNT, workerpool.REQUEST_TIMEOUT, result)

	workerpool.RunWorker(workerPool)
	go workerPool.PushURL(urls)
	go workerPool.GetResult()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	workerPool.Stop()
}
