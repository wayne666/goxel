package goxeler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (g *Goxeler) Run() {
	g.requests = make(chan *request, g.BlockCount)
	//g.results = make(chan *result, g.BlockCount)
	g.runWorkers()
	//close(g.results)

}

func (g *Goxeler) runWorkers() {
	go g.makeRequest() // wait request from g.requests chan

	g.wg.Add(g.BlockCount)
	for i := 0; i < g.BlockCount; i++ {
		rangeStart, rangeEnd := g.calRangeHeader(i)
		// send request to requests chan
		g.requests <- &request{
			blockNum: i,
			retry:    3,
			rangeStartEnd: rangeStartEnd{
				start: rangeStart,
				end:   rangeEnd,
			},
		}
	}
	close(g.requests)
	g.wg.Wait()
}

func (g *Goxeler) makeRequest() {
	for {
		select {
		case request, ok := <-g.requests:
			if !ok {
				fmt.Printf("\n")
				fmt.Println("########## All requests have been sent ##########")
				g.bar = newPb(g.BlockCount)
				return
			}
			go g.downloadFile(request)
		}
	}
}

func (g *Goxeler) downloadFile(request *request) {
	rangeHeader :=
		"bytes=" + strconv.Itoa(request.rangeStartEnd.start) + "-" + strconv.Itoa(request.rangeStartEnd.end)

	req := cloneRequest(g.HttpRequest)
	req.Header.Set("Range", rangeHeader)

	c := http.Client{}
	resp, err := c.Do(req)
	// If HTTP request failed, retry 3 times
	if resp.StatusCode != 206 && err != nil && request.retry-1 > 0 {
		fmt.Printf("Block %d download failed, [error] %v, retrying...\n", err)
		// send request to requests chan, to download again
		g.requests <- request
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
	fmt.Println("Connection ", request.blockNum+1, " has Download.")

	g.wg.Done()
	return
}

func (g *Goxeler) calRangeHeader(blockNum int) (rangeStart, rangeEnd int) {
	rangeStart = blockNum * g.BlockSize
	rangeEnd = rangeStart + g.BlockSize - 1
	if blockNum == (g.BlockCount - 1) {
		rangeEnd = g.FileSize
	}
	return
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
