package goxeler

import (
	"fmt"
	"time"
)

//TODO print some donwload data
func (g *Goxeler) printBar() {

	for {
		select {
		case result := <-g.Results:
			g.Mutex.Lock()
			g.bar.Increment()
			g.Mutex.Unlock()
			println("Connection ", result.blockNum, " has completed.")
			// 减一的原因是blockNum是从0开始的
			if result.blockNum == (g.BlockCount - 1) {
				return
			}
		//case failed := <-g.Fails:
		//	println("Connection ", failed.blockNum, " has failed")
		//	return
		case <-time.After(time.Second * 1):
			//TODO continue print progress bar
			fmt.Println("Ready to download file ...")
		}
	}

}
