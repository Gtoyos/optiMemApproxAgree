# 2-Process Approximate Agreement with Optimal shared memory

This project introduces an algorithm for reaching approximate agreement between two processes using an optimal size of shared memory on each communication round.

Particularly, we use `uint8` as it is the smalles primitive **atomic** type in go.

Particularly we use single-writer multiple-reader (SWMR) shared variables. Particularlt, we consider that the primitive type `uint8` is **atomic**. That is, if a the variable is written by a single process and read at the same time the behaviour is not undefined: it will be either read the previous variable's value or the new updated value.

Note that `uint8` is atomit in any CPU that handles words bigger than 8 bits (current CPU work with 64 bit words).


For calculating the agreement level we calculate the distance in the protocol complex as follows: 

A point x in the protocol complex is named from 0 to x/3^r. Where x ranges from 0 to 3^r. 
Then we take the initial input values to map the input complex to the segment andwe get:

x' = min(a,b) + x/3^r |b-a|

Of course this gives that the difference between consecutives points is:

delta = |b-a| / 3^r 

So if we want the agreement to reach a certain precision Agg. We have at least to do this many rounds:

r > 1/ln(3) ( ln(|b-a|) - ln(Agg) )

With this formula we can map this input to any input in the domain.

### Motivation

My objective was to implement a distribued system algorithm and what best that an interesting one that 
I made myself and is novel :). Also to get some go programming experience. 