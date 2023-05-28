package main

import (
	"fmt"
	"os"
	"github.com/kamil-s-solecki/haze/cliargs"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/mutation"
	"github.com/kamil-s-solecki/haze/report"
	"github.com/kamil-s-solecki/haze/reportable"
)

func readRawRequest(rqPath string) []byte {
	fmt.Println("Request file:", rqPath)
	rawRq, _ := os.ReadFile(rqPath)
	return rawRq
}

func main() {
	args := cliargs.ParseArgs()

	rq := http.Parse(readRawRequest(args.RequestFile))
	addr := args.Host

	reportDir := report.MakeReportDir()
	fmt.Println("Report dir:", reportDir)

	matchers, filters := reportable.FromArgs(args)

	rq.Send(addr)
	for  _, mut := range mutation.Mutate(rq, mutation.AllMutations(), mutation.AllMutatables()) {
		res := mut.Send(addr)
		if reportable.IsReportable(res, matchers, filters) {
			fname := report.Report(mut.Raw(addr), res.Raw, reportDir)
			fmt.Printf("Crash! %s (%s)\n", res, fname)
		}
	}
}
