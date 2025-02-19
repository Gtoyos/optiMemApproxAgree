package main

import (
	"flag"
	"fmt"
	"math"
)

func main() {
	//Flags
	inputValue1 := flag.Float64("input1", 0.0, "Input value of process 0 (compatible with type)")
	inputValue2 := flag.Float64("input2", 1.0, "Input value of process 1 (compatible with type)")
	inputAgreementLvl := flag.Float64("Agreement", 0.000001, "Agreement level. Numeric difference between the two values")

	flag.Parse()

	a := min(*inputValue1, *inputValue2)
	b := max(*inputValue1, *inputValue2)

	//We calculate the required number of rounds to reach the agreement level in function of the distance between values
	r := int(math.Ceil(1 / math.Log(3) * (math.Log(math.Abs(b-a)) - math.Log(*inputAgreementLvl))))
	fmt.Printf("--- Approximate agreement 2-process task ---\n")
	fmt.Printf("Process 0 input: %f\n", a)
	fmt.Printf("Process 1 input: %f\n", b)
	fmt.Printf("Agreement level: %f\n", *inputAgreementLvl)
	fmt.Printf("--------------------------------------------\n")
	fmt.Printf("Agreement level requires %d communication rounds.\n", r)
	//initialize shared snapshots
	snapshots := make([]*SnapshotAtomic[uint8], r)
	for i := 0; i < r; i++ {
		snapshots[i] = NewSnapshotAtomic[uint8](2)
	}
	//channels to return protocol decision
	var c1, c2 chan float64
	c1 = make(chan float64)
	c2 = make(chan float64)

	//start agreement protocol. Generate random delay.
	go func() {
		//wait a random time between 10 and 20 ms
		//time.Sleep(time.Duration(rand.Intn(100-20+1)+10) * time.Millisecond)
		agreement_protocol(0, r, snapshots, c1)
	}()
	go func() {
		//wait a random time between 10 and 20 ms
		//time.Sleep(time.Duration(rand.Intn(100-20+1)+10) * time.Millisecond)
		agreement_protocol(1, r, snapshots, c2)
	}()

	//Print decision values
	//loop to wait for both decisions
	for i := 0; i < 2; i++ {
		select {
		case v := <-c1:
			v = min(a, b) + math.Abs(b-a)*v
			fmt.Printf("Process0 decision output: %f\n", v)
		case v := <-c2:
			v = min(a, b) + math.Abs(b-a)*v
			fmt.Printf("Process1 decision output: %f\n", v)
		}
	}
	fmt.Println("Agreement task compleated.")
}
