package cliargs

import (
	"flag"
	"fmt"
	"strings"
)

type flagName struct {
	long, short string
}

type group struct {
	name      string
	flagNames []flagName
}

var groups []group

func registerFlag(groupName string, fn flagName) {
	for i := range groups {
		g := &groups[i]
		if g.name == groupName {
			g.flagNames = append(g.flagNames, fn)
			return
		}
	}
	groups = append(groups, group{groupName, []flagName{fn}})
}

func printUsage() {
	PrintBanner()
	fmt.Println("USAGE:")
	fmt.Println("  haze [OPTION]... [REQUEST_FILE]...")
	for _, g := range groups {
		fmt.Printf("\n%v:\n", g.name)
		for _, f := range g.flagNames {
			lookup := flag.CommandLine.Lookup(f.long)
			printFlag(f, lookup.Usage, lookup.DefValue)
		}
	}
}

func printFlag(fn flagName, usage, defValue string) {
	ln := "  -" + fn.long
	if fn.short != "" {
		ln += ", -" + fn.short
	}
	ln += strings.Repeat(" ", keyLen-len(ln))
	ln += usage
	if defValue != "" {
		ln += ". (Default: " + defValue + ")"
	}
	fmt.Println(ln)
}
