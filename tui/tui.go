package tui

import (
	"bufio"
	"fmt"
	"github.com/kamil-s-solecki/haze/cliargs"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/progress"
	"log"
	"os"
	"strconv"
	"sync"
)

type Tui struct {
	buff     *bufio.Writer
	mu       sync.Mutex
	errorLog *log.Logger
}

func Create() Tui {
	return Tui{
		buff:     bufio.NewWriter(os.Stdout),
		errorLog: log.New(os.Stdout, "ERROR: ", 0),
	}
}

func (t *Tui) FuzzNewFile(rfile string) {
	t.printf("<< %v >>\n", rfile)
}

func (t *Tui) FuzzNewRequest(rq http.Request) {
	t.printf(" * %v %v\n", rq.Method, rq.RequestUri)
}

func (t *Tui) Crash(res http.Response, fname string) {
	t.printf("(!)  Crash:      %s (%s)\n", res, fname)
}

func (t *Tui) Probe(probe http.Response) {
	t.printf("     Probe:      %v\n", probe)
}

func (t *Tui) EmptyLine() {
	t.printf("\n")
}

func (t *Tui) Fatal(err error) {
	t.errorLog.Fatal(err)
}

func (t *Tui) Error(err error) {
	t.errorLog.Println(err)
}

func (t *Tui) PrintBanner() {
	t.println("               .**.        ")
	t.println("            .. haze ..     ")
	t.println("               `**`        ")
}

func (t *Tui) PrintInfo(args cliargs.Args, reportDir string) {
	entries := []entry{
		{"Target", args.Host},
	}

	if !args.ProbeOnly {
		entries = append(entries, entry{"Report dir", reportDir})
		entries = append(entries, entry{"Threads", strconv.Itoa(args.Threads)})
	}

	if args.Proxy != "" {
		entries = append(entries, entry{"Proxy", args.Proxy})
	}

	t.printTable(entries)
	t.EmptyLine()
}

func (t *Tui) printf(format string, a ...any) {
	defer t.mu.Unlock()
	defer t.buff.Flush()
	t.mu.Lock()

	fmt.Fprintf(t.buff, format, a...)
}

func (t *Tui) println(a ...any) {
	defer t.mu.Unlock()
	defer t.buff.Flush()
	t.mu.Lock()

	fmt.Fprintln(t.buff, a...)
}

func (t *Tui) ProgressBar(total int) progress.Bar {
	return progress.Start(total, t.buff, &t.mu)
}
