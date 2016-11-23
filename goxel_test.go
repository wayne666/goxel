package main

import (
	"testing"
)

func TestBlockSizeCal(t *testing.T) {
	blockSize := blockSizeCal(8, 100)
	if blockSize != 12 {
		t.Errorf("Calculate block size error")
	}
}

func TestGetContentLength(t *testing.T) {
	// The size of license is 19930
	contentLen := getContentLength("http://www.cnperler.com/license.txt")
	if contentLen != 19930 {
		t.Errorf("Get Content Length error")
	}
}
