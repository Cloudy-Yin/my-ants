package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"mq_es_cache/my-ants/ants"
)

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demoFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func main() {
	defer ants.Release()
	var wg sync.WaitGroup
	runTimes := 5

	// // 1. Use the common pool.

	// // syncCalculateSum := func() {
	// // 	demoFunc()
	// // 	wg.Done()
	// // }
	// for i := 0; i < runTimes; i++ {
	// 	wg.Add(1)
	// 	_ = ants.Submit(&ants.Task{
	// 		Handler: func(v ...interface{}) {
	// 			defer wg.Done()
	// 			time.Sleep(time.Second * 2)
	// 			fmt.Println("this is handler func for task:", v[0].(int), v[1].(int), v[2].(int))
	// 		},
	// 		Params: []interface{}{i, i + 2, i + 4},
	// 	})
	// }
	// wg.Wait()
	// fmt.Printf("running goroutines: %d\n", ants.Running())
	// fmt.Printf("finish all tasks.\n")

	// 2. Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(10, func(v ...interface{}) {
		fmt.Println("this is handler func for task:", v)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(i, i+2, i+4)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)

	// 	// Use the MultiPool and set the capacity of the 10 goroutine pools to unlimited.
	// 	// If you use -1 as the pool size parameter, the size will be unlimited.
	// 	// There are two load-balancing algorithms for pools: ants.RoundRobin and ants.LeastTasks.
	// 	mp, _ := ants.NewMultiPool(10, -1, ants.RoundRobin)
	// 	defer mp.ReleaseTimeout(5 * time.Second)
	// 	for i := 0; i < runTimes; i++ {
	// 		wg.Add(1)
	// 		_ = mp.Submit(syncCalculateSum)
	// 	}
	// 	wg.Wait()
	// 	fmt.Printf("running goroutines: %d\n", mp.Running())
	// 	fmt.Printf("finish all tasks.\n")

	// 	// Use the MultiPoolFunc and set the capacity of 10 goroutine pools to (runTimes/10).
	// 	mpf, _ := ants.NewMultiPoolWithFunc(10, runTimes/10, func(i interface{}) {
	// 		myFunc(i)
	// 		wg.Done()
	// 	}, ants.LeastTasks)
	// 	defer mpf.ReleaseTimeout(5 * time.Second)
	// 	for i := 0; i < runTimes; i++ {
	// 		wg.Add(1)
	// 		_ = mpf.Invoke(int32(i))
	// 	}
	// 	wg.Wait()
	// 	fmt.Printf("running goroutines: %d\n", mpf.Running())
	// 	fmt.Printf("finish all tasks, result is %d\n", sum)
	// 	if sum != 499500*2 {
	// 		panic("the final result is wrong!!!")
	// 	}
}
