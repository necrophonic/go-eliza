# go-eliza

[![](https://godoc.org/github.com/necrophonic/go-eliza?status.svg)](http://godoc.org/github.com/necrophonic/go-eliza) [![Go Report Card](https://goreportcard.com/badge/github.com/necrophonic/go-eliza)](https://goreportcard.com/report/github.com/necrophonic/go-eliza)

```
import "github.com/necrophonic/go-eliza"
```

Simple library implementation of Eliza chatbot in Go

Based on the script detailed at [How Eliza Works](http://www.chayden.net/eliza/instructions.txt)

```go
package main

import (
	"bufio"
	"fmt"
	"os"

	eliza "github.com/necrophonic/go-eliza"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Hi, how can I help you?")

	for {
		fmt.Print("> ")
		said, _ := reader.ReadString('\n')

		if said == "bye\n" {
			break
		}

		response, err := eliza.AnalyseString(string(said))
		if err != nil {
			panic(err)
		}
		fmt.Println("[Eliza] " + response)
	}

	fmt.Println("Goodbye")
	os.Exit(0)
}
```

License
-------

Copyright (c) 2016 J Gregory

Released under MIT license, see [LICENSE](LICENSE) for details.
