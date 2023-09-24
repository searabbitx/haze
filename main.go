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
	atui.PrintBanner()
	args := cliargs.ParseArgs()
	http.SetupTransport(args.Proxy)

	reportDir := ""
	if !args.ProbeOnly {
		reportDir = report.MakeReportDir(args.OutputDir)
	}
	atui.PrintInfo(args, reportDir)
	
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

func parseRequestsFromFile(rfile string, args cliargs.Args) (result []http.Request) {
	raw := readRawRequest(rfile)
	if !args.Har {
		result = []http.Request{http.Parse(raw)}
	} else {
		result = http.ParseHar(raw, args.Host)
	}

	if args.Cookies != "" {
		result = overwriteCookies(result, args)
	}

	if len(args.Headers) > 0 {
		result = overwriteHeaders(result, args)
	}

	return
}

func readRawRequest(rqPath string) []byte {
	rawRq, _ := os.ReadFile(rqPath)
	return rawRq
}

func overwriteCookies(rqs []http.Request, args cliargs.Args) []http.Request {
	result := []http.Request{}
	for _, rq := range rqs {
		result = append(result, rq.WithCookieString(args.Cookies))
	}
	return result
}

func overwriteHeaders(rqs []http.Request, args cliargs.Args) []http.Request {
	result := []http.Request{}
	for _, rq := range rqs {
		for _, h := range args.Headers {
			rq = rq.WithHeaderString(h)
		}
		result = append(result, rq)
	}
	return result
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
