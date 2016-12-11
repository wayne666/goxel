package goxeler

import (
	"github.com/cheggaaa/pb"
	"net/http"
	"os"
	"sync"
	"time"
)

type rangeStartEnd struct {
	start int
	end   int
}

type request struct {
	blockNum      int
	retry         int
	rangeStartEnd rangeStartEnd
}

//type result struct {
//	blockNum      int
//	rangeStartEnd rangeStartEnd
//	statusCode    int
//}

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
	bar          *pb.ProgressBar
	requests     chan *request
	successCount int
}

func newPb(size int) (bar *pb.ProgressBar) {
	bar = pb.New(size)
	bar.SetRefreshRate(time.Millisecond)
	bar.ShowSpeed = true
	bar.Start()
	return
}
