package main

import (
	"fmt"
	"os"
	"github.com/kamil-s-solecki/haze/http"
)

func main() {
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
	rq := http.Parse(rawRq)

	fmt.Println(rq)

	rq.Send("http://localhost:9090")
}
