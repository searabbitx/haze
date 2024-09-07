package cliargs

import (
	"flag"
	"fmt"
	"strings"
)

const (
	keyLen = 18
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
	fmt.Println("\nUSAGE:")
	fmt.Println("  haze [OPTION]... [REQUEST_FILE]...")
	fmt.Println("\nARGS:")
	printArg("REQUEST_FILE", []string{
		"File(s) containing the raw http request(s)",
		"in case of .har files pass the -har flag",
		"only the har entries which match the target (-t) value will be fuzzed",
	})
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
	if defValue != "" && defValue != "[  ]" {
		ln += ". (Default: " + defValue + ")"
	}
	fmt.Println(ln)
}

func printArg(name string, usage []string) {
	ln := "  " + name
	ln += strings.Repeat(" ", keyLen-len(ln))
	ln += usage[0]
	for i := 1; i < len(usage); i++ {
		ln += "\n" + strings.Repeat(" ", keyLen) + usage[i]
	}
	fmt.Println(ln)
}
