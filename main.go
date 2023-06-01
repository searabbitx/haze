package main

import (
	"fmt"
	"log"
	"os"
	"github.com/kamil-s-solecki/haze/cliargs"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/mutation"
	"github.com/kamil-s-solecki/haze/progress"
	"github.com/kamil-s-solecki/haze/report"
	"github.com/kamil-s-solecki/haze/reportable"
)

var ErrorLog *log.Logger

func readRawRequest(rqPath string) []byte {
	rawRq, _ := os.ReadFile(rqPath)
	return rawRq
}

func probe(rq http.Request, addr string) {
	probe, err := rq.Send(addr)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	fmt.Println("Probe:            ", probe)
}

func fuzz(args cliargs.Args, rq http.Request, reportDir string) {
	matchers, filters := reportable.FromArgs(args)
	muts := mutation.Mutate(rq, mutation.AllMutations(), mutation.AllMutatables())
	bar := progress.Start(len(muts))
	for  _, mut := range muts {
		res, err := mut.Send(args.Host)
		if err != nil {
			ErrorLog.Println(err)
		}
		if reportable.IsReportable(res, matchers, filters) {
			fname := report.Report(mut.Raw(args.Host), res.Raw, reportDir)
			fmt.Printf("-={*}=- Crash!     %s (%s)\n", res, fname)
		}
		bar.Next()
	}
}

func main() {
	ErrorLog = log.New(os.Stdout, "ERROR: ", 0)
	args := cliargs.ParseArgs()

	rq := http.Parse(readRawRequest(args.RequestFile))
	if args.ProbeOnly {
		cliargs.PrintInfo(args, "")
		probe(rq, args.Host)
		os.Exit(0)
	}

	reportDir := report.MakeReportDir(args.OutputDir)
	cliargs.PrintInfo(args, reportDir)
	probe(rq, args.Host)
	fmt.Println("\n  ...   Fuzzing    ...")

	fuzz(args, rq, reportDir)
}
