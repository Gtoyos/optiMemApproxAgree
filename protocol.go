package main

import (
	"math"
)

type Dir uint8

const (
	Bot Dir = iota
	L
	R
)

// canonical agreement protocol where process 0 has input 0 and process 1 has input 1.
// Processes communicate to decide a value between 0 and 1 that at 1/2^r from each other.
// The algorithm is wait-free.
// Implementation of On the Bit Complexity of iterated memory:
// https://doi.org/10.1007/978-3-031-60603-8_25
func agreement_protocol(input int, rounds int, snapshots []*SnapshotAtomic[uint8], result_ch chan float64) {
	//process id, same as input value (0 or 1)
	pid := input
	//initial state
	s := 0
	for i := 0; i < rounds; i++ {
		//encoding function inside Write
		snapshots[i].Write(uint8(2-(s%2)), pid)
		//read the value of the other process
		v := snapshots[i].Snap()[(1 - pid)]
		//compute the next state. Recall 0->bot, 1->L, 2->R. See pseudocode in paper
		if v == 2 { //Move to the right
			s = (pid-1)*(3*s-1+2*(s%2)) + pid*(3*s+2*(1-(s%2)))
		}
		if v == 0 { //Keep center
			s = 3*s + pid
		}
		if v == 1 { //Move to the left
			s = (pid-1)*(3*s+2*(1-(s%2))) + pid*(3*s-1+2*(s%2))
		}
	}
	result_ch <- float64(2*s+pid) / (math.Pow(3, float64(rounds)))
}
