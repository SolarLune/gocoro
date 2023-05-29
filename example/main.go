package main

import (
	"fmt"
	"time"

	"github.com/solarlune/gocoro"
)

var tickCount int
var progress = -1
var maxProgress = 20

// Here's the coroutine function to run in our Coroutine. It takes an Execution object that is automatically
// created and passed to the coroutine function by Coroutine.Run(). The Execution object allows us to do some
// easy coroutine manipulation (yield, wait until something happens, etc).
// If we want to pause execution, we can call Execution.Yield(). If we want to end execution early,
// just return from the function like usual. If you want to end execution from *outside* this function
// (i.e. with Coroutine.Stop()), then you need to receive that signal by checking Execution.Stopped() from within
// the coroutine function, and return early.
func coroutineFunction(exe *gocoro.Execution) {

	fmt.Printf("\nTick #%d: Let's start the script and wait three seconds.\n", tickCount)

	exe.Wait(time.Second * 3)

	fmt.Printf("\nTick #%d: Excellent! Let's wait 35 ticks this time.\n", tickCount)

	exe.WaitTicks(35)

	fmt.Printf("\nTick #%d: Let's fill this progress bar:\n", tickCount)

	exe.Wait(time.Second * 2)

	fmt.Println("")

	for progress < maxProgress-1 {
		progress++
		exe.Yield()
	}

	progress = -1

	fmt.Printf("\nTick #%d: Excellent, again!\n", tickCount)

	exe.Wait(time.Second)

	fmt.Printf("\nTick #%d: OK, script's over, let's go home!\n", tickCount)

	exe.Wait(time.Second)

}

func main() {

	// Create a new coroutine.
	co := gocoro.NewCoroutine()

	// Run the script. It actually will technically only start when we call Coroutine.Update() below.
	// If we want, we can pass arguments through Coroutine.Run().
	co.Run(coroutineFunction)

	// Coroutine.Running() is thread-safe, as with all of the functions.
	for co.Running() {

		// Update the script. This function call will run the coroutine thread for as long as is necessary,
		// until it either yields (where it will continue the next time Update() is called) or finishes.
		co.Update()

		fmt.Print(".")

		// Draw the progress bar when it's time to do so
		if progress >= 0 {
			pro := "["

			for i := 0; i < maxProgress; i++ {
				if i > progress {
					pro += " "
				} else {
					pro += "█"
				}
			}
			pro += "]"

			fmt.Println(pro)
		}

		tickCount++

		time.Sleep(time.Millisecond * 100)

	}

	fmt.Println("\n\nCoroutine finished!")

	time.Sleep(time.Second * 1)

}
