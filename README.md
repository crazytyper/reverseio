# reverseio
go package for reading files from end to start.

```go
package main

import (
	"fmt"
	"strings"

	"github.com/crazytyper/reverseio"
)

func main() {
	r := reverseio.NewLineReader(strings.NewReader("Hello\nworld!"))

	line, _ := r.ReadString()
	fmt.Println(line)

	line, _ = r.ReadString()
	fmt.Println(line)

	// Output:
	// world!
	// Hello
}
```
