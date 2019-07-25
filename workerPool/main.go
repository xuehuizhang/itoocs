package main

import (
	"fmt"
	"math/rand"
)

type job struct {
	id     int
	number int
}

type result struct {
	job *job
	sum int
}

//计算
func calc(j job, r chan<- *result) {
	var n = j.number
	var sum int
	for n > 0 {
		temp := n % 10
		sum += temp
		n = n / 10
	}
	res := &result{
		job: &j,
		sum: sum,
	}

	r <- res
}

func worker(jobsChan chan *job, resultChan chan *result) {
	for j := range jobsChan {
		calc(*j, resultChan)
	}
}

func workerPool(n int, jobsChan chan *job, resultChan chan *result) {
	for i := 0; i < n; i++ {
		go worker(jobsChan, resultChan)
	}
}

func printResult(res chan *result) {
	for r := range res {
		fmt.Printf("id=%d number=%d  sum=%d\n", r.job.id, r.job.number, r.sum)
	}
}

//worker 池
func main() {
	var jobsChan = make(chan *job, 1000)
	var resultChan = make(chan *result, 1000)

	workerPool(128, jobsChan, resultChan)
	go printResult(resultChan)

	i := 0
	for {
		i++
		num := rand.Int()
		j := &job{
			id:     i,
			number: num,
		}
		jobsChan <- j
	}
}
