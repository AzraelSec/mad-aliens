package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AzraelSec/mad-aliens/pkg/engine"
)

const (
	DEFAULT_MAX_ROUND = 10000
	DEFAULT_ALINES_N  = 10
)

var (
	i = flag.String("i", "", "input file to read world definition from")
	m = flag.Int("m", DEFAULT_MAX_ROUND, "max number of rounds to run")
	n = flag.Int("n", DEFAULT_ALINES_N, "number of aliens to deploy")
)

func init() {
	flag.Parse()
}

func main() {
	// world definition file is required
	if *i == "" {
		flag.Usage()
		return
	}

	file, err := readFile(*i)
	if err != nil {
		log.Fatalf("Impossible to read file %s: %s", *i, err)
	}

	execEngine, err := engine.NewEngine(*n, *m, file)
	if err != nil {
		log.Fatalf("An error occurred during engine initialization: %s", err)
	}

	status, err := execEngine.Run()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Execution completed: %s", engine.ExecStatusString(status))
		log.Printf(
			"| %d survived aliens | %d stuck aliens |",
			execEngine.World.CountAliveAliens(),
			execEngine.World.StuckAliens,
		)
		fmt.Println(execEngine.World)
	}
}

func readFile(p string) (io.Reader, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(f), nil
}
