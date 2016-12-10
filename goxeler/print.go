package goxeler

//TODO print some donwload data
func (g *Goxeler) printProgress() {

	for {
		select {
		case result, ok := <-g.results:
			g.bar.Increment()
			//g.BlockResponseCount++
			println("Connection ", result.blockNum+1, " has completed.")

			//g.wg.Done()
			//return
			if !ok {
				g.bar.FinishPrint("File has download")
			}
		default:
		}
	}
}
