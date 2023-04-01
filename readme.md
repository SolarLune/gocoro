# gocoro 🏃‍♂️ ➡️ 🧍‍♂️

gocoro is a package for basic coroutines for Go. The primary reason I made this was for creating cutscenes for gamedev with Go.

## What's a coroutine?

Normally, a coroutine is just a function that you can pause and resume execution on at will (or possibly even freely jump around at will). Coroutines could have a variety of uses, but are particularly good for scripting cutscenes in games because cutscenes frequently pause and pick up execution (for example, when waiting for some amount of time, displaying a message window, or animating or moving characters to another location).

## Why did you make this package?

Because coroutines are cool and useful, and Go has almost everything necessary for coroutines (like jumping between labels, pausing / blocking execution, etc). Combine that with Go being a fundamentally good, simple, and not too verbose programming language, and it seems like a good solution to this problem, rather than implementing your own scripting language or using a slice of cutscene actions or something.

## How do I use it?

`go get github.com/solarlune/gocoro`

## Example

```go

func script(exe *gocoro.Execution) {

    // Use the Execution object to pause and wait for three seconds.
    exe.Wait(time.Second * 3)

    fmt.Println("Three seconds have elapsed!")

}

func main() {

    // Create a new Coroutine.
    co := gocoro.NewCoroutine()

    // Run the script, which is just an ordinary function pointer that takes
    // an execution object, which is used to control coroutine execution.
    co.Run(script)
    
    for co.Running() {

        // While the coroutine runs, we call Coroutine.Update(). This allows
        // the coroutine to execute, but also gives control back to the main
        // thread when it's yielding so we can do other stuff.
        co.Update()

    }

    // We're done with the coroutine!

}

```

## What's a gocoro.Coroutine?

Internally, a `*gocoro.Coroutine` is just a goroutine running a customizeable function. This means that it executes on another thread. However, gocoro uses a channel to alternately block execution on the calling thread or the coroutine thread, allowing you to pause execution and resume it at will. Because the operating thread alternates between the two, there's no opportunity for race conditions if they both touch the same data.

## Anything else?

Not really, that's it. Peace~