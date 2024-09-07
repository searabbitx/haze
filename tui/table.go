package tui

import (
	"strings"
)

const (
	keyLen = 18
)

type entry struct{ key, val string }

func (t *Tui) printTable(es []entry) {
	max := 0
	lns := []string{}
	for _, e := range es {
		ln := "  " + e.key
		ln += strings.Repeat(" ", keyLen-len(ln))
		values := strings.Split(e.val, "\n")
		ln += ":  " + values[0]
		for _, v := range values[1:] {
			ln += "\n" + strings.Repeat(" ", keyLen) + "   " + v
		}
		lns = append(lns, ln)
		if len(ln) > max {
			max = len(ln)
		}
	}

	bar := strings.Repeat("-", max+2)

	t.println(bar)
	for _, ln := range lns {
		t.println(ln)
	}
	t.println(bar)
}
