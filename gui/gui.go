package gui

import (
	"bufio"
	"fmt"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/progress"
	"os"
	"sync"
)

type Gui struct {
	buff *bufio.Writer
	mu   sync.Mutex
}

func Create() Gui {
	g := Gui{buff: bufio.NewWriter(os.Stdout)}
	return g
}

func (g Gui) FuzzNewFile(rfile string) {
	g.printf("... ( %v ) ...\n", rfile)
}

func (g Gui) FuzzNewRequest(rq http.Request) {
	g.printf("      %v %v\n", rq.Method, rq.RequestUri)
	g.printf("      ---\n")
}

func (g Gui) EmptyLine() {
	g.printf("\n")
}

func (g Gui) printf(format string, a ...any) {
	defer g.mu.Unlock()
	defer g.buff.Flush()
	g.mu.Lock()

	fmt.Fprintf(g.buff, format, a...)
}

func (g Gui) ProgressBar(total int) progress.Bar {
	return progress.Start(total, g.buff, g.mu)
}
