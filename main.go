package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments.")
		os.Exit(1)
	}

	rq := os.Args[1]
	if _, err := os.Stat(rq); err != nil {
		fmt.Println("Cannot read", rq);
		os.Exit(1)
	}

	fmt.Println("Request file:", rq)
}
