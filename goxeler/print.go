package goxeler

import "fmt"

//#type result struct {
//#	start      int
//#	end        int
//#	statusCode int
//#}
func (g *Goxeler) printBar() {

	for {
		select {
		case result := <-g.Results:
			if result.statusCode == 206 {
				g.bar.Increment()
			}
		default:
			return
		}
	}

}
