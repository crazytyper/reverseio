package reverseio

import (
	"fmt"
	"strings"
)

func Example() {
	r := NewLineReader(strings.NewReader("Hello\nworld!"))

	line, _ := r.ReadString()
	fmt.Println(line)

	line, _ = r.ReadString()
	fmt.Println(line)

	// Output:
	// world!
	// Hello
}
