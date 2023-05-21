package main

import (
	"fmt"
	"os"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/mutation"
	"github.com/kamil-s-solecki/haze/reportable"
)

func readRawRequest() []byte {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments.")
		os.Exit(1)
	}

	rqPath := os.Args[1]
	if _, err := os.Stat(rqPath); err != nil {
		fmt.Println("Cannot read", rqPath)
		os.Exit(1)
	}

	fmt.Println("Request file:", rqPath)
	rawRq, _ := os.ReadFile(rqPath)
	return rawRq
}

func main() {
	rq := http.Parse(readRawRequest())
	addr := "http://localhost:9090"

	rq.Send(addr)

	for  _, mut := range mutation.Mutate(rq, mutation.AllMutations(), mutation.AllMutatables()) {
		res := mut.Send(addr)
		if reportable.IsReportable(res) {
			fmt.Println("Found 500!")
		}
	}
}
