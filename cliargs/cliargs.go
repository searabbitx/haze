package cliargs

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Args struct {
	Host          string
	RequestFile   string
	MatchCodes    string
	MatchLengths  string
	FilterCodes   string
	FilterLengths string
	ProbeOnly     bool
}

type Param struct {
	Long, Short, Help string
	Default           interface{}
}

func ParseArgs() Args {
	args := Args{}
	stringVar("GENERAL", &args.Host, Param{Long: "host", Short: "t", Help: "Target host (protocol://hostname:port)"})
	stringVar("GENERAL", &args.RequestFile, Param{Long: "request", Short: "r", Help: "File containing the raw http request"})
	boolVar("GENERAL", &args.ProbeOnly, Param{Long: "probe", Short: "p", Help: "Send the probe request only"})

	stringVar("MATCHERS", &args.MatchCodes, Param{Long: "mc", Default: "500-599", Help: "Comma-separated list of response codes to report"})
	stringVar("MATCHERS", &args.MatchLengths, Param{Long: "ml", Help: "Comma-separated list of response lengths to report"})

	stringVar("FILTERS", &args.FilterCodes, Param{Long: "fc", Help: "Comma-separated list of response codes to not report"})
	stringVar("FILTERS", &args.FilterLengths, Param{Long: "fl", Help: "Comma-separated list of response lengths to not report"})

	flag.Usage = printUsage

	flag.Parse()

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
	validateRequest(args.RequestFile)
	validateRange(args.MatchCodes)
	validateRange(args.MatchLengths)
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

func validateRequest(request string) {
	if request == "" {
		err("The request file (-r, -request) is required")
	}

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
