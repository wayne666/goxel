package goxeler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (g *Goxeler) Run() {
	g.requests = make(chan *request, g.BlockCount)
	g.stopChan = make(chan struct{}, g.BlockCount)
	g.bar = newPb(g.BlockCount)

	g.runWorkers()
}

func (g *Goxeler) runWorkers() {
	g.wg.Add(g.BlockCount)
	for i := 0; i < g.BlockCount; i++ {
		rangeSE := g.calRangeHeader(i)
		req := &request{
			blockNum:      i,
			rangeStartEnd: rangeSE,
		}
		go func() {
			g.makeRequest(req)
			g.wg.Done()
		}()
	}
	g.wg.Wait()
}

func (g *Goxeler) makeRequest(request *request) {
	select {
	case <-g.stopChan:
		return
	default:
		g.downloadFile(request)
	}
}

func (g *Goxeler) downloadFile(request *request) {
	rangeHeader :=
		"bytes=" + strconv.Itoa(request.rangeStartEnd.start) + "-" + strconv.Itoa(request.rangeStartEnd.end)

	req := cloneRequest(g.HttpRequest)
	req.Header.Set("Range", rangeHeader)

	c := http.Client{}
	resp, err := c.Do(req)

	if resp.StatusCode != 206 && err != nil && request.retry-1 > 0 {
		fmt.Printf("Block %d download failed, [error] %v \n", request.blockNum, err)
		return
	}

	body, bodyErr := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if bodyErr != nil {
		fmt.Printf("Read HTTP Body failed [error] %v\n", bodyErr)
		g.wg.Done()
		return
	}

	// seek start position, write body to file
	g.FH.Seek(int64(request.rangeStartEnd.start), 0)
	g.FH.Write([]byte(body))
	g.bar.Increment()
	fmt.Println("Request ", request.blockNum+1, " has Done.")
	g.successCount++
	if g.successCount == g.BlockCount {
		g.bar.FinishPrint("File has download!")
	}

	return
}

// calculate range header start and end
func (g *Goxeler) calRangeHeader(blockNum int) *rangeStartEnd {
	rangeStart := blockNum * g.BlockSize
	rangeEnd := rangeStart + g.BlockSize - 1
	if blockNum == (g.BlockCount - 1) {
		rangeEnd = g.FileSize
	}
	return &rangeStartEnd{
		start: rangeStart,
		end:   rangeEnd,
	}
}

func (g *Goxeler) Stop() {
	for i := 0; i < g.BlockCount; i++ {
		g.stopChan <- struct{}{}
	}
}

// This function comes from https://github.com/rakyll/hey, Thanks for rakyll
// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
