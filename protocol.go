package main

import (
	"math"
	"time"

	"golang.org/x/exp/rand"
)

type Msg uint8

const (
	Bot Msg = iota
	A
	B
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
		snapshots[i].Write(uint8((s%2)+1), pid)
		//read the value of the other process
		v := snapshots[i].Snap()[(1 - pid)]
		//compute the next state. Depending on parity and pid we move to left or right.
		if Msg(v) == A {
			//fmt.Printf("process %d. received A\n", pid)
			s = (1-pid)*(3*s+pid+((s+1)%2)-(s%2)) + pid*(3*s+pid-((s+1)%2)+(s%2))
		}
		if Msg(v) == Bot {
			//fmt.Printf("process %d. received Bot\n", pid)
			s = 3*s + pid
		}
		if Msg(v) == B {
			//fmt.Printf("process %d. received B\n", pid)
			s = (1-pid)*(3*s+pid-((s+1)%2)+(s%2)) + pid*(3*s+pid+((s+1)%2)-(s%2))
		}
		//We add this to simulate delays in communication.
		time.Sleep(time.Duration(rand.Intn(70-20+1)+10) * time.Millisecond)
	}
	result_ch <- float64(2*s+pid) / (math.Pow(3, float64(rounds)))
}
