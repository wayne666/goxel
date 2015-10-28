package goxeler

import (
	"io/ioutil"
	"net/http"
	//"os"
	"fmt"
	"strconv"
	"sync"
)

func (g *Goxeler) Run() {
	g.run()
}

func (g *Goxeler) run() {
	var wg sync.WaitGroup
	wg.Add(g.BlockCount)

	var (
		BlockCount, BlockSize, FileSize int
	)

	BlockCount = g.BlockCount

	headers := make(chan *HeaderRange, BlockCount)
	results := make(chan int, BlockCount)

	for i := 0; i < BlockCount; i++ {
		go func() {
			g.blockDownload(&wg, headers, results)
		}()
	}

	BlockSize = g.BlockSize
	FileSize = g.FileSize

	for i := 0; i < BlockCount; i++ {

		start := i * BlockSize
		end := start + BlockSize - 1

		if i == (BlockSize - 1) {
			end = FileSize
		}
		fmt.Println(start, end)
		headRangeRef := &HeaderRange{
			start: start,
			end:   end,
		}
		headers <- headRangeRef
	}
	close(headers)

	go func() {
		for result := range results {
			println(result)
		}
	}()

	wg.Wait()
}

func (g *Goxeler) blockDownload(wg *sync.WaitGroup, headers chan *HeaderRange, results chan int) {

	for h := range headers {
		defer wg.Done()
		client := &http.Client{}
		req, err := http.NewRequest("GET", g.Url, nil)

		startStr := strconv.Itoa(h.start)
		endStr := strconv.Itoa(h.end)
		headerVal := "bytes=" + startStr + "-" + endStr
		//fmt.Println(headerVal)
		req.Header.Set("Range", headerVal)

		resp, err := client.Do(req)
		g.checkerr(err)

		body, err := ioutil.ReadAll(resp.Body)
		g.checkerr(err)

		//fh := h.fh
		fh := g.FH
		fh.Seek(int64(h.start), 0)
		fh.Write([]byte(body))
		results <- 1
	}
}

func (g *Goxeler) checkerr(e error) {
	if e != nil {
		panic(e)
	}
}
