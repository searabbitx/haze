package main

import (
	"fmt"
	"os"
	"github.com/kamil-s-solecki/haze/cliargs"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/mutation"
	"github.com/kamil-s-solecki/haze/reportable"
	"github.com/kamil-s-solecki/haze/report"
)

func readRawRequest(rqPath string) []byte {
	if _, err := os.Stat(rqPath); err != nil {
		fmt.Println("Cannot read", rqPath)
		os.Exit(1)
	}

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

	rq.Send(addr)
	for  _, mut := range mutation.Mutate(rq, mutation.AllMutations(), mutation.AllMutatables()) {
		res := mut.Send(addr)
		if reportable.IsReportable(res) {
			fmt.Println("Found 500!")
			report.Report(mut.Raw(addr), res.Raw, reportDir)
		}
	}
}
