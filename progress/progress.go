package progress

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
)

type Bar struct {
	curr, total int
	buff        *bufio.Writer
	mu          *sync.Mutex
}

func Start(total int, buff *bufio.Writer, mu *sync.Mutex) Bar {
	b := Bar{curr: 0, total: total, buff: buff, mu: mu}
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

func (b Bar) End() {
	defer b.mu.Unlock()
	b.mu.Lock()
	b.clear()
}

func (b Bar) clear() {
	defer b.buff.Flush()
	spaces := strings.Repeat(" ", len(b.String()))
	fmt.Fprint(b.buff, "\033[30D\033[0K", spaces, "\033[30D")
	fmt.Fprint(b.buff, "\n")
}
