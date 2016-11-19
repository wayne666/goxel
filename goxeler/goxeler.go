package goxeler

import (
	//	"fmt"
	"github.com/cheggaaa/pb"
	"os"
	"sync"
)

//type HeaderRange struct {
//	start int
//	end   int
//}

type result struct {
	sync.Mutex
	start      int
	end        int
	statusCode int
	blockNum   int
}

type failed struct {
	start      int
	end        int
	statusCode int
	blockNum   int
}

type Goxeler struct {
	// numbers of blocks
	N int
	//request header opts
	//Header *HeaderRange
	// result
	Result int
	// request url
	Url string
	// download file size
	FileSize int
	// block count
	BlockCount int
	// each block size
	BlockSize int
	//write file handler
	FH *os.File
	// result struct recieve results
	Results chan *result
	Fails   chan *failed
	// progress bar
	bar *pb.ProgressBar
	//timeout
	timeout chan bool
	//response block count
	BlockResponseCount int
}

func newPb(size int) (bar *pb.ProgressBar) {
	bar = pb.New(size)
	bar.Start()
	return
}
