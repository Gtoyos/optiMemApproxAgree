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
		if Msg(v) == A { //A==1
			s = (1-pid)*(3*s+pid+((s+1)%2)-(s%2)) + pid*(3*s+pid-((s+1)%2)+(s%2))
		}
		if Msg(v) == Bot { //Bot==0
			s = 3*s + pid
		}
		if Msg(v) == B { //B==2
			s = (1-pid)*(3*s+pid-((s+1)%2)+(s%2)) + pid*(3*s+pid+((s+1)%2)-(s%2))
		}
	}
	time.Sleep(time.Duration(rand.Intn(100-20+1)+10) * time.Millisecond)
	result_ch <- float64(2*s+pid) / (math.Pow(3, float64(rounds)))
}
