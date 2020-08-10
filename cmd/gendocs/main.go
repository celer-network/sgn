package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/celer-network/sgn/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory] \n", os.Args[0])
		os.Exit(1)
	}
	path := os.Args[1]

	sgndPath := filepath.Join(path, "sgnd")
	os.RemoveAll(sgndPath)
	os.Mkdir(sgndPath, 0755)
	err := doc.GenMarkdownTree(cmd.GetSgndExecutor().Command, sgndPath)
	if err != nil {
		log.Fatal(err)
	}

	sgncliPath := filepath.Join(path, "sgncli")
	os.RemoveAll(sgncliPath)
	os.Mkdir(sgncliPath, 0755)
	err = doc.GenMarkdownTree(cmd.GetSgncliExecutor().Command, sgncliPath)
	if err != nil {
		log.Fatal(err)
	}

	sgnopsPath := filepath.Join(path, "sgnops")
	os.RemoveAll(sgnopsPath)
	os.Mkdir(sgnopsPath, 0755)
	err = doc.GenMarkdownTree(cmd.GetSgnopsExecutor().Command, sgnopsPath)
	if err != nil {
		log.Fatal(err)
	}
}
