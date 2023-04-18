package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	id int
}

func (j *Job) Run() {
	fmt.Printf("Job %d is running...\n", j.id)
	time.Sleep(time.Second)
	fmt.Printf("Job %d is done.\n", j.id)
}

type Scheduler struct {
	jobs       []*Job
	stopChan   chan bool
	waitGroup  sync.WaitGroup
	maxWorkers int
	workerPool chan bool
}

func NewScheduler(maxWorkers, jobCount int) *Scheduler {
	return &Scheduler{
		jobs:       make([]*Job, jobCount),
		stopChan:   make(chan bool),
		maxWorkers: maxWorkers,
		workerPool: make(chan bool, maxWorkers),
	}
}

func (s *Scheduler) AddJob(job *Job) {
	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Start() {
	for i := 0; i < s.maxWorkers; i++ {
		s.waitGroup.Add(1)
		go s.runWorker()
	}

	for _, job := range s.jobs {
		s.workerPool <- true
		go func(j *Job) {
			defer func() {
				<-s.workerPool
			}()
			j.Run()
		}(job)
	}
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
	s.waitGroup.Wait()
}

func (s *Scheduler) runWorker() {
	defer s.waitGroup.Done()

	for {
		select {
		case <-s.stopChan:
			return
		case <-s.workerPool:
			// take a job from the queue and run it
		}
	}
}

func main() {
	maxWorkers := 3
	scheduler := NewScheduler(maxWorkers, 0)

	for i := 0; i < 10; i++ {
		job := Job{id: i + 1}
		scheduler.AddJob(&job)
	}

	scheduler.Start()
	scheduler.Stop()
}
