package goxeler

//TODO print some donwload data
func (g *Goxeler) printBar() {

	for {
		select {
		case result := <-g.Results:
			g.bar.Increment()
			g.BlockResponseCount++
			println("Connection ", result.blockNum+1, " has completed.")

			if g.BlockResponseCount == g.BlockCount {
				g.bar.FinishPrint("File has download!")
				return
			}
		default:
		}
	}
}
