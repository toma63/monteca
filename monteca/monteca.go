
package main

import ("fmt"
	"flag"
	"github.com/toma63/monteca"
)


// main function
func main() {

	// command line args
	runtime := flag.Int("runtime", 1, "run time in minutes")
	cores := flag.Int("cores", 2, "number of cores")
	flag.Parse()

	// compute percent winnings using Monte Carlo
	pctPlayer, pctBank, totaltrials := monteca.MonteCA(*cores, *runtime)

	fmt.Printf("Total trials using %d cores for %d minutes: %d\n", *cores, *runtime, totaltrials)
	fmt.Printf("Percentage for player: %5.2f\n", pctPlayer)
	fmt.Printf("Percentage for bank: %5.2f\n", pctBank)
}
