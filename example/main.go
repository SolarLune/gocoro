package main

import (
	"fmt"
	"time"

	"github.com/solarlune/gocoro"
)

var gameFrame int
var progress = -1
var maxProgress = 20

// Here's the coroutine function to run in our Coroutine. It takes an execution object that
// allows us to do some easy coroutine manipulation (yield, wait until something happens, etc).
// If we want to pause execution, we can call Execution.Yield(). If we want to end execution early,
// just return from the function like usual. If you want to end execution from *outside* this function
// (i.e. with Coroutine.Kill()), then you can pick up on that through Execution.Killed() and return early.
func coroutineFunction(exe *gocoro.Execution) {

	fmt.Printf("\nFrame #%d: Let's start the script and wait three seconds.\n", gameFrame)

	exe.Wait(time.Second * 3)

	fmt.Printf("\nFrame #%d: Excellent! Let's wait 35 ticks this time.\n", gameFrame)

	exe.WaitTicks(35)

	fmt.Printf("\nFrame #%d: Let's fill this progress bar:\n", gameFrame)

	exe.Wait(time.Second * 2)

	fmt.Println("")

	for progress < maxProgress-1 {
		progress++
		exe.Yield()
	}

	progress = -1

	fmt.Printf("\nFrame #%d: Excellent, again!\n", gameFrame)

	exe.Wait(time.Second)

	fmt.Printf("\nFrame #%d: OK, script's over, let's go home!\n", gameFrame)

	exe.Wait(time.Second)

}

func main() {

	// Create a new coroutine.
	co := gocoro.NewCoroutine()

	// Run the script. It actually will technically only start when we call Coroutine.Update() below.
	co.Run(coroutineFunction)

	// co.Running is thread-safe
	for co.Running() {

		// Update the script. This function call will run the coroutine thread for as long as is necessary,
		// until it either yields or finishes.
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

		gameFrame++

		time.Sleep(time.Millisecond * 100)

	}

	fmt.Println("\n\nCoroutine finished!")

	time.Sleep(time.Second * 1)

}
