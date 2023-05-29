package cliargs

import (
	"fmt"
	"strings"
)

func PrintBanner() {
	fmt.Println("               .**.        ")
	fmt.Println("            .. haze ..     ")
	fmt.Println("               `**`        ")
}

func PrintInfo(args Args, reportDir string) {
	PrintBanner()
	/*
		fmt.Println("-------------------------------------")
		fmt.Println("  Target        : ", args.Host)
		fmt.Println("  Request file  : ", args.RequestFile)
		if reportDir != "" {
			fmt.Println("  Report  dir   : ", reportDir)
		}
		fmt.Println("-------------------------------------\n")
	*/
	entries := []entry{
		{"Target", args.Host},
		{"Request file", args.RequestFile},
	}

	if reportDir != "" {
		entries = append(entries, entry{"Report dir", reportDir})
	}

	printTable(entries)
}

type entry struct{ key, val string }

const (
	keyLen = 16
)

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
