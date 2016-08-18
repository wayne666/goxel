package goxeler

import "fmt"

func (g *Goxeler) printBar() {

	for {
		select {
		case result := <-g.Results:
			//if result.statusCode == 206 {
			//	g.bar.Increment()
			//}
			fmt.Printf("I think should print something, statusCode: %d\n", result.statusCode)
		default:
			return
		}
	}

}
