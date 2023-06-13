package cliargs

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Args struct {
	Host          string
	RequestFiles  []string
	OutputDir     string
	Proxy         string
	Threads       int
	MatchCodes    string
	MatchLengths  string
	MatchString   string
	FilterCodes   string
	FilterLengths string
	FilterString  string
	ProbeOnly     bool
	Har           bool
}

type Param struct {
	Long, Short, Help string
	Default           interface{}
}

func ParseArgs() Args {
	args := Args{}
	stringVar("GENERAL", &args.Host, Param{Long: "host", Short: "t", Help: "Target host (protocol://hostname:port)"})
	boolVar("GENERAL", &args.ProbeOnly, Param{Long: "probe", Short: "p", Help: "Send the probe request only"})
	stringVar("GENERAL", &args.OutputDir, Param{Long: "output", Short: "o", Help: "Directory where the report will be created. (Default: cwd)"})
	intVar("GENERAL", &args.Threads, Param{Long: "threads", Short: "th", Default: 10, Help: "Number of threads to use for fuzzing"})
	stringVar("GENERAL", &args.Proxy, Param{Long: "proxy", Short: "x", Help: "Proxy address"})
	boolVar("GENERAL", &args.Har, Param{Long: "har", Help: "Indicate that the request files are in the har format"})

	stringVar("MATCHERS", &args.MatchCodes, Param{Long: "mc", Default: "500-599", Help: "Comma-separated list of response codes to report"})
	stringVar("MATCHERS", &args.MatchLengths, Param{Long: "ml", Help: "Comma-separated list of response lengths to report"})
	stringVar("MATCHERS", &args.MatchString, Param{Long: "ms", Help: "A string to match in response"})

	stringVar("FILTERS", &args.FilterCodes, Param{Long: "fc", Help: "Comma-separated list of response codes to not report"})
	stringVar("FILTERS", &args.FilterLengths, Param{Long: "fl", Help: "Comma-separated list of response lengths to not report"})
	stringVar("FILTERS", &args.FilterString, Param{Long: "fs", Help: "A string to filter in response"})

	flag.Usage = printUsage

	flag.Parse()
	args.RequestFiles = flag.Args()

	validate(args)

	fixArgs(&args)
	return args
}

func stringVar(group string, pvar *string, param Param) {
	registerFlag(group, flagName{param.Long, param.Short})
	deflt := ""
	if param.Default != nil {
		deflt = param.Default.(string)
	}
	flag.StringVar(pvar, param.Long, deflt, param.Help)
	if param.Short != "" {
		flag.StringVar(pvar, param.Short, deflt, "")
	}
}

func intVar(group string, pvar *int, param Param) {
	registerFlag(group, flagName{param.Long, param.Short})
	deflt := 0
	if param.Default != nil {
		deflt = param.Default.(int)
	}
	flag.IntVar(pvar, param.Long, deflt, param.Help)
	if param.Short != "" {
		flag.IntVar(pvar, param.Short, deflt, "")
	}
}

func boolVar(group string, pvar *bool, param Param) {
	registerFlag(group, flagName{param.Long, param.Short})
	deflt := false
	if param.Default != nil {
		deflt = param.Default.(bool)
	}
	flag.BoolVar(pvar, param.Long, deflt, param.Help)
	if param.Short != "" {
		flag.BoolVar(pvar, param.Short, deflt, "")
	}
}

func validate(args Args) {
	validateHost(args.Host)
	validateProxy(args.Proxy)
	validateRequests(args.RequestFiles)
	validateRange(args.MatchCodes)
	validateRange(args.MatchLengths)
	validateOutput(args.OutputDir)
}

func validateHost(host string) {
	if host == "" {
		err("The target host (-t, -host) is required")
	}

	r, _ := regexp.Compile("^https?://([-a-zA-Z.]{1,256})(:[0-9]{1,5})?/?$")
	if !r.MatchString(host) {
		err("The target host should be in format: protocol://hostname:port")
	}
}

func validateProxy(proxy string) {
	if proxy == "" {
		return
	}

	r, _ := regexp.Compile("^(https?|socks[0-9]?)://([-a-zA-Z.]{1,256})(:[0-9]{1,5})?/?$")
	if !r.MatchString(proxy) {
		err("The proxy string should be in format: protocol://hostname:port")
	}
}

func validateRequests(rqs []string) {
	if len(rqs) == 0 {
		err("The request file(s) is required")
	}

	for _, rq := range rqs {
		validateRequest(rq)
	}
}

func validateRequest(request string) {
	fi, e := os.Stat(request)
	if e != nil {
		err("Cannot read: " + request)
	}
	if fi.IsDir() {
		err(request + " is a directory. Please provide a file")
	}
}

func validateRange(val string) {
	if val == "" {
		return
	}

	r, _ := regexp.Compile("^[0-9]+(-[0-9]+)?(,[0-9]+(-[0-9]+)?)*$")
	if !r.MatchString(val) {
		err(fmt.Sprintf("Invalid range: '%v'. Example correct value: '100,200-300,400'", val))
	}
}

func validateOutput(output string) {
	if output == "" {
		return
	}

	fi, e := os.Stat(output)
	if e != nil {
		err("Cannot open: " + output)
	}
	if !fi.IsDir() {
		err(output + " is not a directory. Please provide a directory")
	}
}

func err(msg string) {
	fmt.Println(msg + "\n")
	flag.Usage()
	os.Exit(1)
}

func fixArgs(args *Args) {
	if args.Host[len(args.Host)-1:] == "/" {
		args.Host = args.Host[:len(args.Host)-1]
	}
}
