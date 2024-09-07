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
}

type Param struct {
	Long, Short, Default, Help string
}

func ParseArgs() Args {
	args := Args{}
	stringVar(&args.Host, Param{Long: "host", Short: "t", Help: "Target host (protocol://hostname:port)"})
	stringVar(&args.RequestFile, Param{Long: "request", Short: "r", Help: "File containing the raw http request"})

	stringVar(&args.MatchCodes, Param{Long: "mc", Default: "500-599", Help: "Comma-separated list of response codes to report"})
	stringVar(&args.MatchLengths, Param{Long: "ml", Help: "Comma-separated list of response lengths to report"})

	stringVar(&args.FilterCodes, Param{Long: "fc", Help: "Comma-separated list of response codes to not report"})
	stringVar(&args.FilterLengths, Param{Long: "fl", Help: "Comma-separated list of response lengths to not report"})

	configUsage()

	flag.Parse()

	validate(args)

	fixArgs(&args)
	return args
}

func stringVar(pvar *string, param Param) {
	flag.StringVar(pvar, param.Long, param.Default, param.Help)
	if param.Short != "" {
		flag.StringVar(pvar, param.Short, param.Default, "-"+param.Long+" (shorthand)")
	}
}

func configUsage() {
	flag.Usage = func() {
		PrintBanner()
		fmt.Println("OPTIONS:\n")
		flag.PrintDefaults()
	}
}

func PrintBanner() {
	fmt.Println("               .**.        ")
	fmt.Println("            .. haze ..     ")
	fmt.Println("               `**`        ")
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
