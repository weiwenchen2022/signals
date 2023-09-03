## signals

Simple wrapper for os/signal.

## Install

```sh
go get github.com/weiwenchen2022/signals
```

## Example

```go
package main

import (
	"os"
	"syscall"

	"github.com/weiwenchen2022/signals"
)

func main() {
	stop := signals.Notify(func(os.Signal) bool {
		// reload config
		// true indicates continue serving for further signals.
		return true
	}, syscall.SIGUSR2)
	defer stop()

	select {}
}
```

## Doc
GoDoc: [https://godoc.org/github.com/weiwenchen2022/signals](https://godoc.org/github.com/weiwenchen2022/signals)
