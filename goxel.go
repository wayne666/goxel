package main

import (
	"flag"
	"fmt"
	"github.com/WayneZhouChina/goxel/goxeler"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
)

const (
	headerRegexp = `^([\w-]+):\s*(.+)`
)

type headerSlice []string

func (h *headerSlice) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headerSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}

var (
	headerslice headerSlice
	version     = "0.1.0"
	n           = flag.Int("n", 8, "")
	o           = flag.String("o", "", "")
	verbose     = flag.Bool("v", false, "")
	V           = flag.Bool("V", false, "version")
	h           = flag.Bool("h", false, "help information")
	cpus        = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
)

var usage = ` Usage: goxel [options...] <url>

Options:
	-n  Specify number of connections.
	-H  Custom HTTP header. You can specify as many as the header you needd.
		For example, -H "Accept: text/html" -H "Content-Type: application/xml" .
	-o  Specify local output file.
	-h  Help information.
	-v  More status information.
	-V  Version
	-cpus Number of used cpu cores(Default is current machine cores).
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}

	flag.Var(&headerslice, "H", "")
	flag.Parse()

	n := *n
	output := *o
	verbose := *verbose
	V := *V
	h := *h
	cpus := *cpus

	runtime.GOMAXPROCS(cpus)

	if verbose {
		usageAndExit("")
		return
	}
	if V {
		fmt.Println("goxel version ", version)
		return
	}
	if h {
		usageAndExit("")
		return
	}

	if n <= 0 {
		usageAndExit("n cannot be smaller than 1.")
	}

	if output == "" {
		usageAndExit("o must be assgined the output file.")
	}

	if flag.NArg() < 1 {
		usageAndExit("")
	}

	url := flag.Args()[0]

	// Set HTTP Header
	header := make(http.Header)
	for _, h := range headerslice {
		match, err := parseInputWithRegexp(h, headerRegexp)
		if err != nil {
			usageAndExit(err.Error())
		}
		header.Set(match[0], match[1])
	}

	// Init http request, and add custom header
	req, err := http.NewRequest("GET", url, nil)
	req.Header = header

	// Create output file handler
	fh, err := os.Create(output)
	checkerr(err)
	defer fh.Close()

	// Get download file size, and caculate block count
	blockCount := n
	filesize := getContentLength(url)
	blockSize := blockSizeCalculate(blockCount, filesize)

	(&goxeler.Goxeler{
		HttpRequest: req,
		FileSize:    filesize,
		BlockCount:  blockCount,
		BlockSize:   blockSize,
		Url:         url,
		FH:          fh,
	}).Run()
}

func parseInputWithRegexp(input, regexpString string) ([]string, error) {
	re := regexp.MustCompile(regexpString)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return nil, fmt.Errorf("could not parse header string")
	}
	return matches, nil
}

func getContentLength(url string) int {
	req, err := http.Head(url)
	checkerr(err)

	filesize := req.Header.Get("Content-Length")
	if filesize == "" {
		fmt.Println("file size invalid")
		os.Exit(0)
	}

	filesizeInt, _ := strconv.Atoi(filesize)
	return filesizeInt
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

func blockSizeCalculate(blockCount, filesize int) (blockSize int) {
	blockSize = (filesize - filesize%blockCount) / blockCount
	return
}

func checkerr(e error) {
	if e != nil {
		panic(e)
	}
}
