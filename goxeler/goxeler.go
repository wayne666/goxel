package goxeler

import (
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
)

type rangeStartEnd struct {
	start int
	end   int
}

type request struct {
	blockNum      int
	retry         int
	rangeStartEnd *rangeStartEnd
}

type Goxeler struct {
	wg          sync.WaitGroup
	HttpRequest *http.Request
	Url         string
	// download file size
	FileSize int
	// block count
	BlockCount int
	BlockSize  int
	//open write filehandle
	FH *os.File
	//progress bar
	start        time.Time
	bar          *pb.ProgressBar
	requests     chan *request
	successCount int
	stopChan     chan struct{}
	ProxyAddr    *url.URL
	results      chan *result
}

type result struct {
	err           error
	statusCode    int
	duration      time.Duration
	connDuration  time.Duration // connection setup(DNS lookup + Dial up) duration
	dnsDuration   time.Duration // dns lookup duration
	reqDuration   time.Duration // request "write" duration
	resDuration   time.Duration // response "read" duration
	delayDuration time.Duration // delay between response and request
	contentLength int64
}

func newPb(size int) (bar *pb.ProgressBar) {
	bar = pb.New(size)
	bar.SetRefreshRate(time.Millisecond)
	bar.ShowSpeed = true
	bar.Start()
	return
}
