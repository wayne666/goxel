package goxeler

import (
	"os"
)

type headerRange struct {
	start int
	end   int
	//	url   string
	fh *os.File
}

type Goxeler struct {
	// numbers of blocks
	N int
	// request header opts
	header *headerRange
	// result
	result int
	// request url
	url string
}
