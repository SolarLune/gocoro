10/25/23: Development on gocoro has been discontinued to be superceded by another repo for scripting event sequences for game development, [routine](https://github.com/solarlune/routine). This is due to a fundamental problem that cannot easily be solved in function execution (which is that functions cannot easily be exited from external to the coroutine). While my personal development and testing revealed that `runtime.Goexit()` could be used to quit running coroutines while they're running, in practice for game development (at least when using Ebitengine), accessing data from various goroutines (even if they're not running at the same time) may still lead to timing issues. This might be because Ebitengine modifies data internally from another goroutine. 

It might be possible to bypass the issues I've run into, but for now, I think it best to pursue for a less esoteric solution for simply and less "magically" scripting events.

# gocoro 🏃‍♂️ ➡️ 🧍

gocoro is a package for basic coroutines for Go. The primary reason I made this was for creating cutscenes for gamedev with Go.

## What's a coroutine?

Normally, a coroutine is just a function that you can pause and resume execution on at will (or possibly even freely jump around at will). Coroutines could have a variety of uses, but are particularly good for scripting cutscenes in games because cutscenes frequently pause and pick up execution (for example, when a cutscene, say, waits a certain amount of time, the game isn't frozen, but rather continues playing audio, taking input, and updating the screen. The coroutine is running and waiting, but the game continues).

## What's a gocoro.Coroutine?

Internally, a `*gocoro.Coroutine` is just a goroutine running a customizeable function. This means that it probably would execute on another thread (as to the Go runtime's determinations). However, gocoro uses a channel to alternately block execution on the calling thread or the coroutine thread, allowing you to pause execution and resume it at will. Because only one thread out of the two is active at any given time, there's no opportunity for race conditions if they touch the same data; this means that gocoro's Coroutines should be inherently thread-safe.

## Why did you make this package?

Because coroutines are cool and useful, and Go has almost everything necessary for coroutines (like jumping between labels, pausing / blocking execution, etc). Combine that with Go being a fundamentally good, simple, and not too verbose programming language, and it seems like a good solution to this problem, rather than implementing your own scripting language or using a slice of cutscene actions or something.

## How do I use it?

`go get github.com/solarlune/gocoro`

## Example

```go

func script(exe *gocoro.Execution) {

    // Use the Execution object to yield execution and wait for three seconds.
    exe.YieldTime(time.Second * 3)

    fmt.Println("Three seconds have elapsed!")

}

func main() {

    // Create a new Coroutine.
    co := gocoro.NewCoroutine()

    // Run the script, which is just an ordinary function pointer that takes
    // an execution object, which is automatically provided when a Coroutine is 
    // running and is used to control coroutine execution.
    // You can pass extra arguments to the function through the Run command as well.
    co.Run(script)
    
    for co.Running() {

        // While the coroutine runs, we call Coroutine.Update(). This allows
        // the coroutine to execute, but also gives control back to the main
        // thread when it's yielding so we can do other stuff, like take input
        // or update a game's screen.
        co.Update()

    }

    // After Running() is over, we're done with the coroutine! Of course, you can just check
    // this with an if statement, rather than a for statement.

}

```

## Anything else?

Not really, that's it. Peace~