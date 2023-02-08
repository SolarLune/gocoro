package gocoro

import (
	"errors"
	"sync/atomic"
	"time"
)

// Coroutine represents a coroutine that executes alternately with the main / calling
// thread.
type Coroutine struct {
	routine   func(*Execution)
	running   *atomic.Bool
	yield     chan bool
	execute   chan bool
	Execution *Execution
}

// NewCoroutine creates and returns a new Coroutine instance.
func NewCoroutine() Coroutine {
	co := Coroutine{
		yield:   make(chan bool),
		execute: make(chan bool),
		running: &atomic.Bool{},
	}
	co.Execution = &Execution{coroutine: &co}
	return co
}

// Run runs the given coroutine function. Note that the function takes the Coroutine
// as an argument to allow for a variety of methods of defining this function (as a
// literal in the Run() function call, as a pointer to a pre-defined function, etc).
// Run will return an error if the coroutine is already running.
func (co *Coroutine) Run(coroutineFunc func(exe *Execution)) error {

	if !co.running.Load() {

		co.running.Store(true)

		co.routine = coroutineFunc

		go func() {
			// Send something on execute first so the script doesn't update until we
			// call Coroutine.Update() the first time.
			co.execute <- true
			co.routine(co.Execution)
			wasRunning := co.running.Load()
			co.running.Store(false)
			// If the coroutine wasn't running anymore, then we shouldn't push anything to yield to unblock the coroutine at the end
			if wasRunning {
				co.yield <- true
			}
		}()

		return nil

	} else {
		return errors.New("Coroutine is already running")
	}

}

// Running returns whether the Coroutine is running or not.
func (co *Coroutine) Running() bool {
	return co.running.Load()
}

// Update waits for the Coroutine to pause, either as a yield or when the Coroutine is finished. If the
// Coroutine isn't running anymore, Update doesn't do anything.
func (co *Coroutine) Update() {
	if co.running.Load() {
		<-co.execute // Wait to pull from the execute channel, indicating the coroutine can run
		<-co.yield   // Wait to pull from the yield channel, indicating the coroutine has paused / finished
	}
}

// Stop stops running the Coroutine and allows the CoroutineExecution to pick up on this fact to end gracefully.
// Note that this does not kill the coroutine, which internally runs in a goroutine - you'll need to detect this and
// end the coroutine from the coroutine function yourself.
func (co *Coroutine) Stop() {
	wasRunning := co.running.Load()
	co.running.Store(false)
	if wasRunning {
		<-co.execute // Pull from the execute channel so the coroutine can get out of the yield and realize it's borked
	}
}

var ErrorCoroutineStopped = errors.New("Coroutine requested to be stopped")

// Execution represents a means to easily and simply manipulate coroutine execution from your running coroutine function.
type Execution struct {
	coroutine *Coroutine
}

// Yield yields execution in the coroutine function, allowing the main / calling thread to continue.
// The coroutine will pick up from this point when Coroutine.Update() is called again.
// If the Coroutine has exited already, then this will immediately return with ErrorCoroutineStopped.
func (exe *Execution) Yield() error {

	if !exe.coroutine.Running() {
		return ErrorCoroutineStopped
	}

	exe.coroutine.yield <- true   // Yield; we're done
	exe.coroutine.execute <- true // Put something in the execute channel when we're ready to pick back up if we're not done

	return nil

}

// Stopped returns true if the coroutine was requested to be stopped through Coroutine.Stop(). You can check this in your
// coroutine to exit early and clean up the coroutine as desired.
func (exe *Execution) Stopped() bool {
	return !exe.coroutine.Running()
}

// Wait waits the specified duration time, yielding execution in the Coroutine if the time has yet to elapse.
// Note that this function only checks the time in increments of however long the calling thread takes between calling Coroutine.Update().
// So, for example, if Coroutine.Update() is run, say, once every 20 milliseconds, then that's the fidelity of your waiting duration.
// If the Coroutine has stopped prematurely, then this will immediately return with ErrorCoroutineStopped.
func (exe *Execution) Wait(duration time.Duration) error {
	start := time.Now()
	for {

		if time.Since(start) >= duration {
			return nil
		} else {
			if err := exe.Yield(); err != nil {
				return err
			}
		}
	}
}

// WaitTicks waits the specified number of ticks, yielding execution if the number of ticks have yet to elapse. A tick is defined by one instance
// of Coroutine.Update() being called.
// If the Coroutine has stopped prematurely, then this will immediately return with ErrorCoroutineStopped.
func (exe *Execution) WaitTicks(tickCount int) error {
	for {

		if tickCount == 0 {
			return nil
		} else {
			tickCount--
			if err := exe.Yield(); err != nil {
				return err
			}
		}

	}

}

// WaitUntil pauses the Coroutine until the provided Completer's Done() function returns true.
// If the Coroutine has stopped prematurely, then this will immediately return with ErrorCoroutineStopped.
func (exe *Execution) WaitUntil(completer Completer) error {

	for {

		if completer.Done() {
			return nil
		} else {
			if err := exe.Yield(); err != nil {
				return err
			}
		}
	}

}

// Do pauses the running Coroutine until the provided function returns true.
// If the Coroutine has stopped prematurely, then this will immediately return with ErrorCoroutineStopped.
func (exe *Execution) Do(doFunc func() bool) error {

	for {
		if doFunc() {
			return nil
		} else {
			if err := exe.Yield(); err != nil {
				return err
			}
		}
	}

}

// Completer provides an interface of an object that can be used to pause a Coroutine until it is completed.
// If the Completer's Done() function returns true, then the Coroutine will advance.
type Completer interface {
	Done() bool
}
