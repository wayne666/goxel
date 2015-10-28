package main

import (
	"./goxeler"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

type HeaderRange struct {
	start int
	end   int
	url   string
	fh    *os.File
}

type RespOpts struct {
	start  int
	end    int
	status int
}

type Goxeler struct {
	// numbers of blocks
	N int
	// result
	result int
	// download file size
	fileSize int
	// block count
	blockCount int
}

var (
	// block count
	n      = flag.Int("n", 8, "")
	header = flag.String("H", "", "")
	ua     = flag.String("U", "", "")
	o      = flag.String("o", "", "")

	verbose = 0
)

var usage = ` Usage: goxel [options...] <url>

Options:
	-n  Numbers of blocks to run.
	-H  Add header string.
	-U  Set user agent.
	-v  More status information.
	-o  Specify local output file.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}
	flag.Parse()

	n := *n
	output := *o

	if n <= 0 {
		usageAndExit("n cannot be smaller than 1.")
	}

	if output == "" {
		usageAndExit("o must be assgined the output file.")
	}

	if flag.NArg() < 1 {
		usageAndExit("")
	}

	var (
		blockSize, fileSize, blockCount int
		url                             string
	)

	url = flag.Args()[0]
	blockSize = n * 1024 * 1024

	fileSize = fileSizeCal(url)
	blockCount = blockCountCal(blockSize, fileSize)

	fh, err := os.Create(output)
	checkerr(err)
	defer fh.Close()

	(&goxeler.Goxeler{
		N:          n,
		FileSize:   fileSize,
		BlockCount: blockCount,
		BlockSize:  blockSize,
		Url:        url,
		FH:         fh,
	}).Run()

}

func analyUrl(url string) string {
	return url
}

func fileSizeCal(url string) int {
	resp, err := http.Head(url)
	checkerr(err)

	fileSize := resp.Header.Get("Content-Length")

	if fileSize == "" {
		fmt.Println("file size invalid")
		os.Exit(0)
	}

	fileSizeInt, err := strconv.Atoi(fileSize)
	return fileSizeInt
}

func usageAndExit(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func blockCountCal(blockSize, fileSize int) int {

	var blockCount int
	remainder := fileSize % blockSize
	blockCount = fileSize / blockSize

	if remainder != 0 {
		blockCount += 1
	}

	return blockCount
}

func checkerr(e error) {
	if e != nil {
		panic(e)
	}
}
