package goxeler

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
