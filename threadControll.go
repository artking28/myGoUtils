package myGoUtils

import "sync"

type ThreadControl struct {
	channel       chan int
	wg            *sync.WaitGroup
	mu            *sync.Mutex
	activeThreads int
}

// NewThreadControl initializes a new ThreadControl instance with a specified maximum number of threads.
func NewThreadControl(maxThreads int) *ThreadControl {
	ret := ThreadControl{
		wg:            &sync.WaitGroup{},
		mu:            &sync.Mutex{},
		activeThreads: 0,
	}
	if maxThreads > 0 {
		ret.channel = make(chan int, maxThreads)
	}
	return &ret
}

// Lock acquires the mutex lock to protect shared resources.
func (this *ThreadControl) Lock() {
	this.mu.Lock()
}

// Unlock releases the mutex lock.
func (this *ThreadControl) Unlock() {
	this.mu.Unlock()
}

// Begin starts a new thread, defaulting the delta to 1.
func (this *ThreadControl) Begin() {
	this.BeginN(1)
}

// BeginN starts a specified number of new threads, indicated by deltaz
func (this *ThreadControl) BeginN(delta int) {
	this.channel <- 0
	this.mu.Lock()
	this.activeThreads += delta
	this.mu.Unlock()
	this.wg.Add(delta)
}

// Done signals that a thread has completed its work.
func (this *ThreadControl) Done() {
	this.wg.Done()
	this.mu.Lock()
	this.activeThreads--
	this.mu.Unlock()
	<-this.channel // Libera o slot no canal
}

// Wait blocks until all active threads have finished executing.
func (this *ThreadControl) Wait() {
	this.wg.Wait()
}
