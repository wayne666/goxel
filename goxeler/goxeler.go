package goxeler

import (
	//	"fmt"
	"os"
)

type HeaderRange struct {
	start int
	end   int
	//url   string
	//fh    *os.File
}

type Goxeler struct {
	// numbers of blocks
	N int
	//request header opts
	Header *HeaderRange
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
}

type Result struct {
	Start      int
	End        int
	statusCode int
}
