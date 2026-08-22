package main

import (
	"fmt"
	"os"

	"github.com/suer/ssmsh/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
