package bufpool_test

import (
	"fmt"

	"github.com/tomcz/gotools/bufpool"
)

func Example() {
	pool := bufpool.New()

	buf := pool.Borrow()
	defer pool.Return(buf)

	fmt.Fprint(buf, "Hello")
	fmt.Println(buf.String())

	// Output:
	// Hello
}
