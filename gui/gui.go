package gui

import (
	"bufio"
	"fmt"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/progress"
	"log"
	"os"
	"sync"
)

type Gui struct {
	buff     *bufio.Writer
	mu       sync.Mutex
	errorLog *log.Logger
}

func Create() Gui {
	return Gui{
		buff:     bufio.NewWriter(os.Stdout),
		errorLog: log.New(os.Stdout, "ERROR: ", 0),
	}
}

func (g Gui) FuzzNewFile(rfile string) {
	g.printf("... ( %v ) ...\n", rfile)
}

func (g Gui) FuzzNewRequest(rq http.Request) {
	g.printf("      %v %v\n", rq.Method, rq.RequestUri)
	g.printf("      ---\n")
}

func (g Gui) Crash(res http.Response, fname string) {
	g.printf("    ! Crash:       %s (%s)\n", res, fname)
}

func (g Gui) Probe(probe http.Response) {
	g.printf("      Probe:      %v\n", probe)
}

func (g Gui) EmptyLine() {
	g.printf("\n")
}

func (g Gui) Fatal(err error) {
	g.errorLog.Fatal(err)
}

func (g Gui) Error(err error) {
	g.errorLog.Println(err)
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
