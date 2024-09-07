package cliargs

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	keyLen = 18
)

func PrintBanner() {
	fmt.Println("               .**.        ")
	fmt.Println("            .. haze ..     ")
	fmt.Println("               `**`        ")
}

func PrintInfo(args Args, reportDir string) {
	PrintBanner()
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

	printTable(entries)
}

type entry struct{ key, val string }

func printTable(es []entry) {
	max := 0
	lns := []string{}
	for _, e := range es {
		ln := "  " + e.key
		ln += strings.Repeat(" ", keyLen-len(ln))
		ln += ":  " + e.val
		lns = append(lns, ln)
		if len(ln) > max {
			max = len(ln)
		}
	}

	bar := strings.Repeat("-", max+2)

	fmt.Println(bar)
	for _, ln := range lns {
		fmt.Println(ln)
	}
	fmt.Println(bar)
	fmt.Println("")
}
