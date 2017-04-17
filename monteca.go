
package monteca

import ("math/rand"
	"time"
	//"fmt"
)

const (
	MAXBET float64 = 100.0
)

type winnings struct {
	bank float64
	player float64
}
	

// initialize a random caBoard
// coin flip to decide to bet on any square
// value between min and max for bet amount
// board represented as a float64 slice with 6 locations
func makeRandCABoard() []float64 {

	// initialize random number generator
	now := time.Now()
	r := rand.New(rand.NewSource(now.UnixNano()))
	
	// slice of 6 floats
	cboard := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

	for i := 0 ; i < 6 ; i++ {

		// bet on a given square based on a coin flip
		if r.Uint32() & 0x00000001 == 1 {
			cboard[i] = r.Float64() * MAXBET
			//fmt.Printf("bet for %d is %5.2f\n", i, cboard[i])
		} else {
			//fmt.Printf("no bet for %d\n", i)
			cboard[i] = 0.0
		}
	}
	return cboard
} 

// generate one die roll - a number between 0 and 5
func rollOne() int32 {

	// initialize random number generator
	now := time.Now()
	r := rand.New(rand.NewSource(now.UnixNano()))

	return r.Int31n(6) // 0 - 5

}

// individual worker
// compute the amount of bet money retained by player and bank
// accumulate results for runtime minutes
func caWorker(resultCh chan winnings, trialsCh chan int64, runtime int) {

	result := winnings{0.0, 0.0} // bank, player
	now := time.Now()
	totalTrials := int64(0)

	// loop accumulating the result until time is up
	for time.Since(now) < (time.Duration(runtime) * time.Minute) {
		totalTrials += 1
		board := makeRandCABoard()
		hits := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
		
		// roll three dice
		for i := 0 ; i < 3 ; i++ {
			hits[rollOne()] += 1.0
		}

		for i := 0 ; i < 6 ; i++ {
			if hits[i] > 0.99 {
				won := hits[i] * board[i]
				result.player += won + board[i]
				result.bank -= won
			} else {
				result.bank += board[i]
			}
		}
		
	}
	resultCh <- result
	trialsCh <- totalTrials
}


// compute winning percentage for the game of crown and anchor
//   cores specifies the number of parallel workers
//   runtime specifies the run time in minutes
func MonteCA(cores int, runtime int) (pctPlayer float64, pctBank float64, totalTrials int64) {

	totalCh := make(chan winnings, cores)
	trialsCh := make(chan int64, cores)

	// launch the workers
	for i := 0 ; i < cores ; i++ {
		go caWorker(totalCh, trialsCh, runtime)
	}
	
	// drain the channels
	totalWinnings := winnings{0.0, 0.0}
	totalTrials = 0
	var total winnings

	for i := 0 ; i < cores ; i++ {
		total = <- totalCh
		totalWinnings.bank += total.bank
		totalWinnings.player += total.player
		totalTrials +=  <- trialsCh
	}

	playerPlusBank := totalWinnings.bank + totalWinnings.player
	pctPlayer = (totalWinnings.player / playerPlusBank) * 100
	pctBank = 100.0 - pctPlayer
	
	return 
}


