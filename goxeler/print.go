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
			if result.blockNum == (g.BlockCount - 1) {
				return
			}
		case <-time.After(time.Second * 1):
			//TODO continue print progress bar
			fmt.Println("....")
		}
	}

}
