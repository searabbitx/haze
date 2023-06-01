package progress

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type Bar struct {
	curr, total int
	buff        *bufio.Writer
	mu          sync.Mutex
}

func Start(total int) Bar {
	b := Bar{curr: 0, total: total, buff: bufio.NewWriter(os.Stdout)}
	return b
}

func (b *Bar) Next() {
	defer b.mu.Unlock()
	b.mu.Lock()
	b.curr++
	b.update()
}

func (b Bar) Log(msg string) {
	defer b.mu.Unlock()
	b.mu.Lock()
	fmt.Fprint(b.buff, msg, "\n")
	b.buff.Flush()
}

func (b Bar) update() {
	defer b.buff.Flush()
	fmt.Fprint(b.buff, "\033[30D\033[0K", b, "\033[30D")
}

const spinChars = `|/-\`

func (b Bar) spinner() byte {
	return spinChars[b.curr%len(spinChars)]
}

func (b Bar) String() string {
	return fmt.Sprintf("     %c [ %v / %v ] %c", b.spinner(), b.curr, b.total, b.spinner())
}
