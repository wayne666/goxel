package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type headRange struct {
	start int
	end   int
	url   string
	fh    *os.File
}

var (
	header = flag.String("H", "", "")
	ua     = flag.String("U", "", "")
	o      = flag.String("o", "", "")

	bc = flag.Int("bc", 8, "")
	//bs = flag.Int("bs", 2, "")

	verbose = 0
)

var usage = ` Usage: goxel [options...] <url>

Options:
	-bn Numbers of blocks to run.
	-H  Add header string.
	-U  Set user agent.
	-v  More status information.
	-o  Specify local output file.
`

func main() {
	var url string
	url = "http://112.253.22.162/9/m/y/d/z/mydziciabmgopabgscdoihgvrnzdvd/he.yinyuetai.com/F2A5015080D4156E413EE6227BB0E8C2.flv?sc=e01f47f6135c5339&br=3101&vid=2398390&aid=38959&area=KR&vst=0&ptp=mv&rd=yinyuetai.com"

	bc := *bc

	var (
		blockSize, fileSizeInt, blockCount int
	)

	blockSize = bc * 1024 * 1024

	fileSizeInt = fileSizeCal(url)
	blockCount = blockCountCal(blockSize, fileSizeInt)

	// create dest file handle
	fh, err := os.Create("./goxel_download.flv")
	checkerr(err)

	defer fh.Close()

	// add wg wait group
	var wg sync.WaitGroup

	wg.Add(blockCount)

	jobs := make(chan *headRange, blockCount)
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
			end = fileSizeInt
		}
		//fmt.Println(start, end)
		headRangeRef := &headRange{
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

func blockCountCal(blockSize int, fileSize int) int {

	remainder := fileSize % blockSize
	blockCount = fileSize / blockSize

	if remainder != 0 {
		blockCount += 1
	}

	return blockCount
}

func segChanDownload(wg *sync.WaitGroup, ch chan *headRange, result chan int) {
	//func segChanDownload(wg *sync.WaitGroup, ch chan string, result chan int) {

	for val := range ch {
		defer wg.Done()
		url := val.url
		start := val.start
		end := val.end
		f := val.fh

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

	//fmt.Println("[start] ", start)
	//f.Seek(int64(start), 0)
	//checkerr(err)
	//f.Write([]byte(body))

	//defer f.Close()
	//checkerr(err)
	return
}

func checkerr(e error) {
	if e != nil {
		panic(e)
	}
}
