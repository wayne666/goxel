package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	//"reflect"
)

type headerRange struct {
	start int
	end   int
	url   string
	fh    *os.File
}

type respOpts struct {
	start  int
	end    int
	status int
}

type goxeler struct {
	// numbers of blocks
	N int
	// request header opts
	header *headerRange
	// result
	result int
}

var (
	n      = flag.Int("n", 8, "") // block count
	header = flag.String("H", "", "")
	ua     = flag.String("U", "", "")
	o      = flag.String("o", "", "")

	//bc = flag.Int("bc", 8, "")
	//bs = flag.Int("bs", 2, "")

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
		//blockSize, fileSize, blockCount int
		blockSize, fileSize, blockCount int
		url                             string
	)

	url = flag.Args()[0]
	blockSize = n * 1024 * 1024

	fileSize = fileSizeCal(url)
	blockCount = blockCountCal(blockSize, fileSize)

	fmt.Println(fileSize)
	fmt.Println(blockCount)

	fh, err := os.Create(output)
	checkerr(err)
	defer fh.Close()

	// add wg wait group
	var wg sync.WaitGroup

	wg.Add(blockCount)

	jobs := make(chan *headerRange, blockCount)
	//jobs := make(chan string, segCount)
	result := make(chan int, blockCount)

	for i := 0; i < blockCount; i++ {
		go func() {
			segChanDownload(&wg, jobs, result)
		}()
	}

	for i := 0; i < blockCount; i++ {

		start := i * blockSize
		end := start + blockSize - 1
		println(end)

		if i == (blockSize - 1) {
			end = fileSize
		}
		//fmt.Println(start, end)
		headRangeRef := &headerRange{
			start: start,
			end:   end,
			fh:    fh,
			url:   url,
		}
		jobs <- headRangeRef
	}
	close(jobs)

	//for i := 0; i < segCount; i++ {
	//	println(<-result)
	//}
	go func() {
		for val := range result {
			println(val)
		}
	}()

	wg.Wait()
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

//func main() {
//	var url string
//	url = "http://112.253.22.162/9/m/y/d/z/mydziciabmgopabgscdoihgvrnzdvd/he.yinyuetai.com/F2A5015080D4156E413EE6227BB0E8C2.flv?sc=e01f47f6135c5339&br=3101&vid=2398390&aid=38959&area=KR&vst=0&ptp=mv&rd=yinyuetai.com"
//
//	bc := *bc
//
//
//	blockSize = n * 1024 * 1024
//
//	fileSizeInt = fileSizeCal(url)
//	blockCount = blockCountCal(blockSize, fileSizeInt)
//
//	// create dest file handle
//	fh, err := os.Create("./goxel_download.flv")
//	checkerr(err)
//
//	defer fh.Close()
//
//	// add wg wait group
//	var wg sync.WaitGroup
//
//	wg.Add(blockCount)
//
//	jobs := make(chan *headRange, blockCount)
//	//jobs := make(chan string, segCount)
//	result := make(chan int, blockCount)
//
//	for i := 0; i < blockCount; i++ {
//		go func() {
//			segChanDownload(&wg, jobs, result)
//		}()
//	}
//
//	for i := 0; i < blockCount; i++ {
//
//		start := i * blockSize
//		end := start + blockSize - 1
//		println(end)
//
//		if i == (blockSize - 1) {
//			end = fileSizeInt
//		}
//		//fmt.Println(start, end)
//		headRangeRef := &headRange{
//			start: start,
//			end:   end,
//			fh:    fh,
//			url:   url,
//		}
//		jobs <- headRangeRef
//	}
//	close(jobs)
//
//	//for i := 0; i < segCount; i++ {
//	//	println(<-result)
//	//}
//	go func() {
//		for val := range result {
//			println(val)
//		}
//	}()
//
//	wg.Wait()
//}
//

func blockCountCal(blockSize, fileSize int) int {

	var blockCount int
	remainder := fileSize % blockSize
	blockCount = fileSize / blockSize

	if remainder != 0 {
		blockCount += 1
	}

	return blockCount
}

//
func segChanDownload(wg *sync.WaitGroup, ch chan *headerRange, result chan int) {

	defer wg.Done()
	var (
		start int
		end   int
		url   string
		f     *os.File
	)
	for val := range ch {
		url = val.url
		start = val.start
		end = val.end
		f = val.fh

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		startStr := strconv.Itoa(start)
		endStr := strconv.Itoa(end)
		headerVal := "bytes=" + startStr + "-" + endStr
		fmt.Println(headerVal)
		req.Header.Set("Range", headerVal)

		resp, err := client.Do(req)
		checkerr(err)

		for key, value := range resp.Header {
			fmt.Printf("%s: %s\n", key, value)
		}
		//println(resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		//ioutil.WriteFile("/tmp/test.flv", body, 0644)

		////fmt.Println("[start] ", start)
		f.Seek(int64(start), 0)
		//checkerr(err)

		f.Write([]byte(body))
		//f.Write([]byte("hello"))
		result <- 1
	}

	//url := hr.url
	//start := hr.start
	//end := hr.end
	//client := &http.Client{}
	//req, err := http.NewRequest("GET", url, nil)
	//startStr := strconv.Itoa(start)
	//endStr := strconv.Itoa(end)

	//headerVal := "bytes=" + startStr + "-" + endStr
	//println(headerVal)
	//req.Header.Set("Range", headerVal)

	//resp, err := client.Do(req)
	//checkerr(err)

	//for key, value := range resp.Header {
	//	fmt.Printf("%s: %s\n", key, value)
	//}
	////println(resp.StatusCode)
	//body, err := ioutil.ReadAll(resp.Body)
	//defer resp.Body.Close()
	////ioutil.WriteFile("/tmp/test.flv", body, 0644)

	////fmt.Println("[start] ", start)
	////f.Seek(int64(start), 0)
	////checkerr(err)
	////f.Write([]byte(body))

	////defer f.Close()
	////checkerr(err)
	return
}

func checkerr(e error) {
	if e != nil {
		panic(e)
	}
}
