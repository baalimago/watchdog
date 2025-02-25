# Watchdog

Super simple watchdog implementation based on the existing timer functionality.

## Example

This is an example implementation where the objective is to gather buffer data as long as it comes in during a minimum update frequence, and then flush the buffer.

```
ctx, cancel := context.WithCancel(ctx.Background())
defer cancel()

var buffer []interface{}
m := sync.Mutex{}
flushData := func() {
    m.Lock()
    defer m.Unlock()
    flush(buffer)
}

// Create a new watchdog and enter callback to call (such as a data flush operation)
w := watchdog.New(5 * time.Millisecond, flushData)

go func() {
BREAK:
    for {
        select {
            case <- ctx.Done():
                break OUTER
            case bufferableData := <- sourceChan:
                // Pet the watchdog to keep it from flushing for another 5 milliseconds
                w.Pet()
                m.Lock()
                buffer := append(buffer, bufferableData)
                m.Unlock()
        }
    }
}()

```

To avoid infinite Pets, use the `w.Stop()` functionality along with some separate timeout functionality (such as `context.WithTimeout`).
