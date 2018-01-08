package goxeler

import (
	"fmt"
	"time"
)

type report struct {
	avgTotal float64
	fastest  float64
	slowest  float64
	average  float64
	rps      float64

	avgConn  float64
	avgDNS   float64
	avgReq   float64
	avgRes   float64
	avgDelay float64

	results chan *result
	done    chan bool
	total   time.Duration

	errorDist      map[string]int
	statusCodeDist map[int]int
	sizeTotal      int64
}

func newReport(results chan *result) *report {
	return &report{
		errorDist:      make(map[string]int),
		done:           make(chan bool, 1),
		statusCodeDist: make(map[int]int),
		results:        results,
	}
}

func runReporter(r *report) {
	// Loop will continue until channel is closed
	for res := range r.results {
		if res.err != nil {
			r.errorDist[res.err.Error()]++
		} else {
			r.avgTotal += res.duration.Seconds()
			r.avgConn += res.connDuration.Seconds()
			r.avgDelay += res.delayDuration.Seconds()
			r.avgDNS += res.dnsDuration.Seconds()
			r.avgReq += res.reqDuration.Seconds()
			r.avgRes += res.resDuration.Seconds()

			r.statusCodeDist[res.statusCode]++
			if res.contentLength > 0 {
				r.sizeTotal += res.contentLength
			}
		}
	}
	r.done <- true
}

func (r *report) finalize(total time.Duration) {
	r.total = total
	r.print()
}

func (r *report) print() {
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total:\t%4.4f secs\n", r.total.Seconds())
	if r.sizeTotal > 0 {
		fmt.Printf("  Total data:\t%d bytes\n", r.sizeTotal)
	}
	fmt.Println()
	fmt.Println("Details (average):")
	fmt.Printf("  DNS+dialup %4.4f secs\n", r.avgConn)
	fmt.Printf("  DNS-lookup %4.4f secs\n", r.avgDNS)
	fmt.Printf("  req write  %4.4f secs\n", r.avgReq)
	fmt.Printf("  resp wait  %4.4f secs\n", r.avgDelay)
	fmt.Printf("  resp read  %4.4f secs\n", r.avgRes)
	r.printStatusCodes()

	if len(r.errorDist) > 0 {
		r.printErrors()
	}
	fmt.Printf("\n")
}

// printStatusCodes prints status code distribution.
func (r *report) printStatusCodes() {
	fmt.Printf("\n\nStatus code distribution:\n")
	for code, num := range r.statusCodeDist {
		fmt.Printf("  [%d]\t%d responses\n", code, num)
	}
}

func (r *report) printErrors() {
	fmt.Printf("\nError distribution:\n")
	for err, num := range r.errorDist {
		fmt.Printf("  [%d]\t%s\n", num, err)
	}
}
