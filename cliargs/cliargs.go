package cliargs

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

type Args struct {
	Host        string
	RequestFile string
}

type Param struct {
	Long, Short, Default, Help string
}

func ParseArgs() Args {
	args := Args{}
	stringVar(&args.Host, Param{Long: "host", Short: "t", Help: "Target host (protocol://hostname:port)"})
	stringVar(&args.RequestFile, Param{Long: "request", Short: "r", Help: "File containing the raw http request"})

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
		printBanner()
		fmt.Println("OPTIONS:\n")
		flag.PrintDefaults()
	}
}

func printBanner() {
	fmt.Println("... Haze ...\n")
}

func validate(args Args) {
	validateHost(args.Host)
	validateRequest(args.RequestFile)
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
