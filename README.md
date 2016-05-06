# reverseio
go package for reading files from end to start.

```go
package reverseio

import (
	"fmt"
	"strings"
)

func main() {
	r := NewLineReader(strings.NewReader("Hello\nworld!"))

	line, _ := r.ReadString()
	fmt.Println(line)

	line, _ = r.ReadString()
	fmt.Println(line)

	// Output:
	// world!
	// Hello
}
```
