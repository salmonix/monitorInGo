package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gmon/watch/process"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	register(process.NewWatchedProcess(os.Getpid(), 0))
	enterLoop()
}

func register(p *process.WatchedProcess) {
	body, _ := json.Marshal(p)
	toPost := bytes.NewBuffer(body)
	url := "http://localhost:8080/process" // TODO: do proper url
	mime := "text/json"
	resp, _ := http.Post(url, mime, toPost)
	fmt.Println(resp)
}

func enterLoop() {
	mb := 1024 * 1024 // 1 MB in bytes
	for {
		for s := 8; s < 4*8; s = s + 8 {
			memoryEater := make([]float32, mb*s)
			if cap(memoryEater) > mb*s {
				fmt.Println("impossible")
			}
			loopMessage(fmt.Sprintf(" adding memory eater of %d Mb", s))
			time.Sleep(3000 * time.Millisecond)
		}
		loopMessage("Heavy math")
		primes()
	}
}

func loopMessage(s string) {
	fmt.Printf("%s: %s\n", strconv.FormatInt(time.Now().Unix(), 10), s)
}

func calculateIntensively() {
	timeout := time.After(10 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		case <-timeout:
			return
		case <-tick:
			primes()
		}
	}
}

// copy-paste of https://github.com/Agis-/gofool/blob/master/atkin.go
// the Sieve of Atkin algo
func primes() {
	const N = 1000000
	var x, y, n int
	nsqrt := math.Sqrt(N)
	isPrime := [N]bool{}

	for x = 1; float64(x) <= nsqrt; x++ {
		for y = 1; float64(y) <= nsqrt; y++ {
			n = 4*(x*x) + y*y
			if n <= N && (n%12 == 1 || n%12 == 5) {
				isPrime[n] = !isPrime[n]
			}
			n = 3*(x*x) + y*y
			if n <= N && n%12 == 7 {
				isPrime[n] = !isPrime[n]
			}
			n = 3*(x*x) - y*y
			if x > y && n <= N && n%12 == 11 {
				isPrime[n] = !isPrime[n]
			}
		}
	}

	for n = 5; float64(n) <= nsqrt; n++ {
		if isPrime[n] {
			for y = n * n; y < N; y += n * n {
				isPrime[y] = false
			}
		}
	}

	isPrime[2] = true
	isPrime[3] = true

	primes := make([]int, 0, 1270606)
	for x = 0; x < len(isPrime)-1; x++ {
		if isPrime[x] {
			primes = append(primes, x)
		}
	}
}