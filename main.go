package main

import (
	"os"
	"github.com/kamil-s-solecki/haze/cliargs"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/mutable"
	"github.com/kamil-s-solecki/haze/mutation"
	"github.com/kamil-s-solecki/haze/report"
	"github.com/kamil-s-solecki/haze/reportable"
	"github.com/kamil-s-solecki/haze/workerpool"
	"github.com/kamil-s-solecki/haze/tui"
)

var atui tui.Tui

func main() {
	atui = tui.Create()
	args := cliargs.ParseArgs()
	http.SetupTransport(args.Proxy)

	reportDir := ""
	if !args.ProbeOnly {
		reportDir = report.MakeReportDir(args.OutputDir)
	}
	cliargs.PrintInfo(args, reportDir)
	
	for _, rfile := range args.RequestFiles {
		atui.FuzzNewFile(rfile)
		for _, rq := range parseRequestsFromFile(rfile, args) {
			atui.FuzzNewRequest(rq)
			probe(rq, args.Host)
			if args.ProbeOnly {
				atui.EmptyLine()
			} else {
				fuzz(args, rq, reportDir)
			}
		}
	}
}

func parseRequestsFromFile(rfile string, args cliargs.Args) []http.Request {
	raw := readRawRequest(rfile)
	if !args.Har {
		return []http.Request{http.Parse(raw)}
	}
	return http.ParseHar(raw, args.Host)
}

func readRawRequest(rqPath string) []byte {
	rawRq, _ := os.ReadFile(rqPath)
	return rawRq
}

func probe(rq http.Request, addr string) {
	probe, err := rq.Send(addr)
	if err != nil {
		atui.Fatal(err)
	}
	atui.Probe(probe)
}

func fuzz(args cliargs.Args, rq http.Request, reportDir string) {
	matchers, filters := reportable.FromArgs(args)
	muts := mutation.Mutate(rq, mutation.AllMutations(), mutable.AllMutatables())
	bar := atui.ProgressBar(len(muts))
	pool := workerpool.NewPool(args.Threads)

	for  _, mut := range muts {
		mut := mut
		task := func() {
			res, err := mut.Send(args.Host)
			if err != nil {
				atui.Error(err)
			}
			if reportable.IsReportable(res, matchers, filters) {
				fname := report.Report(mut.Raw(args.Host), res.Raw, reportDir)
				atui.Crash(res, fname)
			}
			bar.Next()
		}
		pool.RunTask(task)
	}
	pool.Wait()
	bar.End()
}
