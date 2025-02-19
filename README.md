# 2-Process Approximate Agreement with Optimal Shared Memory


This project introduces an distributed algorithm for reaching approximate agreement between two processes using an optimal size of shared memory on each communication round.


This implementation is based on the optimal bounded IIS two-process approximate agreement algorithm described in the paper by G. Toyos-Marfurt and P. Kuznetsov, “On the Bit Complexity of Iterated Memory”, available [here](https://arxiv.org/abs/2402.12484).

## About the implementation

We use `uint8` as it is the smalles primitive **atomic** type in go.

We use single-writer multiple-reader (SWMR) shared variables, assuming that `uint8` operations are **atomic**. Specifically, when a uint8 variable is written by one process and read simultaneously, the behavior is deterministic: the reader observes either the previous value or the newly written value.

Note that `uint8` is atomic in any CPU that handles words bigger than 8 bits (current CPUs typically work with 64 bit words).

For demonstration purposes, we provide parallel implementations of the algorithm:
- Go: Using native multithreading with goroutines and channels.
- Rust: Leveraging the [Tokio](https://tokio.rs) asynchronous runtime library.

Code implementations are available in the `/go` and `/rust` directories respectively.

## Agreement precision and rounds

For calculating the agreement, we map the possible values on the protocol complex as follows:

A point $x$ in the protocol complex is defined as:

$$
x = \frac{i}{3^r} \ \ i\in\{0,1,3^r\}
$$

For generic inputs $a$ and $b$, we map the $0$-$1$ output complex to the extended segment as:

$$
x' = \min(a,b) + x \times |b-a|
$$

Therefore, the difference between consecutive points is:

$$
\delta = \frac{|b-a|}{3^r}
$$

To achieve a target agreement precision $\delta$ , the number of rounds  $r$  must satisfy:

$$
r > \frac{1}{\log{3}}(\log{|b-a|} - \log{\delta})
$$
