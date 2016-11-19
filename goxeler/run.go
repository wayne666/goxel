package goxeler

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

func (g *Goxeler) Run() {
	//收集结果result的channel
	g.Results = make(chan *result, g.BlockCount)

	// new 一个 progress bar 显示进度
	g.bar = newPb(g.BlockCount)
	g.run()
	close(g.Results)
	g.FH.Close()
}

func (g *Goxeler) run() {
	go g.printBar()

	var wg sync.WaitGroup
	wg.Add(g.BlockCount)
	for i := 0; i < g.BlockCount; i++ {
		go func(i int) {
			g.blockDownload(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func (g *Goxeler) blockDownload(blockNum int) {
	//拼装 range 的开始和结束
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

	// 拼装 http range header, 正确返回响应码 206
	startStr := strconv.Itoa(rangeStart)
	endStr := strconv.Itoa(rangeEnd)
	headerStr := "bytes=" + startStr + "-" + endStr
	req.Header.Set("Range", headerStr)

	resp, err := client.Do(req)
	g.checkerr(err)

	//获取响应码，正常为206, 因为是range 请求
	statusCode := resp.StatusCode

	body, bodyErr := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// 如果响应有错误
	//if bodyErr != nil {
	//	g.Fails <- &failed{
	//		start:      rangeStart,
	//		end:        rangeEnd,
	//		statusCode: statusCode,
	//		blockNum:   blockNum,
	//	}
	//	return
	//}

	// 如果正确返回，块写入文件，并且把结果给 result chan
	if statusCode == 206 && bodyErr == nil && g.bar != nil {
		g.FH.Seek(int64(rangeStart), 0)
		g.FH.Write([]byte(body))

		g.Results <- &result{
			start:      rangeStart,
			end:        rangeEnd,
			statusCode: statusCode,
			blockNum:   blockNum,
		}
	}
}

func (g *Goxeler) checkerr(e error) {
	if e != nil {
		panic(e)
	}
}
