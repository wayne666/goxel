package goxeler

import (
	//"github.com/cheggaaa/pb"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	//"time"
	"fmt"
)

func (g *Goxeler) Run() {
	g.Results = make(chan *result, g.BlockCount)
	g.bar = newPb(g.BlockCount)
	g.run()
	g.bar.FinishPrint("File has download!")
	close(g.Results)
}

func (g *Goxeler) run() {
	var wg sync.WaitGroup
	blockNumChan := make(chan int, g.BlockCount)

	wg.Add(g.BlockCount)
	for i := 0; i < g.BlockCount; i++ {
		blockNumChan <- i
		go func() {
			g.blockDownload(blockNumChan)
			wg.Done()
		}()
	}
	//g.printBar()
	wg.Wait()
}

func (g *Goxeler) blockDownload(blockNumChan chan int) {
	blockNum := <-blockNumChan

	rangeStart := blockNum * g.BlockSize
	rangeEnd := rangeStart + g.BlockSize - 1
	if blockNum == (g.BlockCount - 1) {
		rangeEnd = g.FileSize
	}
	g.makeRequest(rangeStart, rangeEnd, blockNum)
}

func (g *Goxeler) makeRequest(rangeStart, rangeEnd, blockNum int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", g.Url, nil)

	startStr := strconv.Itoa(rangeStart)
	endStr := strconv.Itoa(rangeEnd)
	fmt.Printf("start: %d, end: %d\n", rangeStart, rangeEnd)
	headerStr := "bytes=" + startStr + "-" + endStr
	req.Header.Set("Range", headerStr)

	resp, err := client.Do(req)
	g.checkerr(err)

	statusCode := 0
	if err == nil {
		statusCode = resp.StatusCode
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if statusCode == 206 && bodyErr == nil && g.bar != nil {
		g.bar.Increment()
	}
	g.checkerr(bodyErr)

	g.FH.Seek(int64(rangeStart), 0)
	g.FH.Write([]byte(body))

	g.Results <- &result{
		start:      rangeStart,
		end:        rangeEnd,
		statusCode: statusCode,
		blockNum:   blockNum,
	}
}

func (g *Goxeler) checkerr(e error) {
	if e != nil {
		panic(e)
	}
}
