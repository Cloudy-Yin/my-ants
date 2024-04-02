package ants

import (
	"runtime/debug"
	"time"
)

// goWorkerWithFunc is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type goWorkerWithFunc struct {
	// pool who owns this worker.
	pool *PoolWithFunc

	// args is a job should be done.
	args chan []interface{}

	// lastUsed will be updated when putting a worker back into queue.
	lastUsed time.Time
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *goWorkerWithFunc) run() {
	w.pool.addRunning(1)
	go func() {
		defer func() {
			w.pool.addRunning(-1)
			w.pool.workerCache.Put(w)
			if p := recover(); p != nil {
				if ph := w.pool.options.PanicHandler; ph != nil {
					ph(p)
				} else {
					w.pool.options.Logger.Printf("worker exits from panic: %v\n%s\n", p, debug.Stack())
				}
			}
			// Call Signal() here in case there are goroutines waiting for available workers.
			w.pool.cond.Signal()
		}()

		for args := range w.args {
			if args == nil {
				return
			}
			w.pool.poolFunc(args...)
			if ok := w.pool.revertWorker(w); !ok {
				return
			}
		}
	}()
}

func (w *goWorkerWithFunc) finish() {
	w.args <- nil
}

func (w *goWorkerWithFunc) lastUsedTime() time.Time {
	return w.lastUsed
}

func (w *goWorkerWithFunc) inputFunc(*Task) {
	panic("unreachable")
}

func (w *goWorkerWithFunc) inputParam(arg ...interface{}) {
	w.args <- arg
}
