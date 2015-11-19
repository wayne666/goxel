package goxeler

import (
	//"github.com/cheggaaa/pb"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	//	"time"
)

func (g *Goxeler) Run() {
	g.Results = make(chan *result, g.BlockCount)
	//g.bar = newPb(g.BlockCount)
	g.run()
	//close(g.Results)
}

func (g *Goxeler) run() {
	var wg sync.WaitGroup
	wg.Add(g.BlockCount)

	var (
		BlockCount, BlockSize, FileSize int
	)

	BlockCount = g.BlockCount

	headers := make(chan *HeaderRange, BlockCount)

	for i := 0; i < BlockCount; i++ {
		go func() {
			g.blockDownload(&wg, headers)
		}()
	}

	BlockSize = g.BlockSize
	FileSize = g.FileSize

	for i := 0; i < BlockCount; i++ {
		start := i * BlockSize
		end := start + BlockSize - 1

		if i == (BlockCount - 1) {
			end = FileSize
			println("=============")
			println(end)
			println("=============")
		}

		headers <- &HeaderRange{
			start: start,
			end:   end,
		}
	}
	close(headers)

	go func() {
		for result := range g.Results {
			//println("==================")
			println(result.start)
		}
	}()
	//for {
	//	aresult := <-g.Results
	//	println(aresult.start)
	//	//if result.statusCode == 206 {
	//	//	g.bar.Increment()
	//	//}
	//	//	}
	//}
	//}()

	wg.Wait()
}

func (g *Goxeler) blockDownload(wg *sync.WaitGroup, headers chan *HeaderRange) {

	fh := g.FH
	for h := range headers {
		defer wg.Done()
		client := &http.Client{}
		req, err := http.NewRequest("GET", g.Url, nil)

		startStr := strconv.Itoa(h.start)
		endStr := strconv.Itoa(h.end)
		headerStr := "bytes=" + startStr + "-" + endStr
		println(headerStr)
		req.Header.Set("Range", headerStr)

		resp, err := client.Do(req)

		code := 0
		if err == nil {
			code = resp.StatusCode
		}
		println(code)

		body, err := ioutil.ReadAll(resp.Body)
		g.checkerr(err)

		fh.Seek(int64(h.start), 0)
		fh.Write([]byte(body))

		g.Results <- &result{
			start:      h.start,
			end:        h.end,
			statusCode: code,
		}
	}
}

func (g *Goxeler) checkerr(e error) {
	if e != nil {
		panic(e)
	}
}
