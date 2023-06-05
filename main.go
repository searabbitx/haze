package main

import (
	"fmt"
	"log"
	"os"
	"github.com/kamil-s-solecki/haze/cliargs"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/mutable"
	"github.com/kamil-s-solecki/haze/mutation"
	"github.com/kamil-s-solecki/haze/progress"
	"github.com/kamil-s-solecki/haze/report"
	"github.com/kamil-s-solecki/haze/reportable"
	"github.com/kamil-s-solecki/haze/workerpool"
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
	fmt.Println("      Probe:      ", probe)
}

func fuzz(args cliargs.Args, rq http.Request, reportDir string) {
	matchers, filters := reportable.FromArgs(args)
	muts := mutation.Mutate(rq, mutation.AllMutations(), mutable.AllMutatables())
	bar := progress.Start(len(muts))
	pool := workerpool.NewPool(args.Threads)

	for  _, mut := range muts {
		mut := mut
		task := func() {
			res, err := mut.Send(args.Host)
			if err != nil {
				ErrorLog.Println(err)
			}
			if reportable.IsReportable(res, matchers, filters) {
				fname := report.Report(mut.Raw(args.Host), res.Raw, reportDir)
				bar.Log(fmt.Sprintf("     !Crash:       %s (%s)\n", res, fname))
			}
			bar.Next()
		}
		pool.RunTask(task)
	}
	pool.Wait()
}

func main() {
	ErrorLog = log.New(os.Stdout, "ERROR: ", 0)
	args := cliargs.ParseArgs()
	http.SetupTransport(args.Proxy)

	reportDir := ""
	if !args.ProbeOnly {
		reportDir = report.MakeReportDir(args.OutputDir)
	}
	cliargs.PrintInfo(args, reportDir)
	
	for _, rfile := range args.RequestFiles {
		fmt.Printf("... ( %v ) ...\n", rfile)
		rq := http.Parse(readRawRequest(rfile))
		if args.ProbeOnly {
			probe(rq, args.Host)
			fmt.Println()
			continue
		}
		probe(rq, args.Host)
		fuzz(args, rq, reportDir)
	}
}
