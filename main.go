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
	runTimes := 200

	// 1. Use the common pool.

	// syncCalculateSum := func() {
	// 	demoFunc()
	// 	wg.Done()
	// }

	resCh := make(chan int, 200)
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = ants.Submit(&ants.Task{
			Handler: func(v ...interface{}) {
				defer wg.Done()
				//time.Sleep(time.Second * 2)
				fmt.Println("this is handler func for task:", v[0].(int), v[1].(int), v[2].(int))
				resCh <- v[0].(int) + v[1].(int) + v[2].(int)
			},
			Params: []interface{}{i, i + 2, i + 4},
		})
	}

	wg.Wait()
	close(resCh)

	sum := 0
	for v := range resCh {
		fmt.Println("resulte:", v)
		sum += v
	}

	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks, sum : %v.\n", sum)

	// // 2. Use the pool with a function,
	// // set 10 to the capacity of goroutine pool and 1 second for expired duration.
	// resCh := make(chan int, 200)
	// p, _ := ants.NewPoolWithFunc(10, func(v ...interface{}) {
	// 	params := v[0].([]interface{})
	// 	fmt.Printf("var name type: %T, value: %v\n", params, params)
	// 	fmt.Println("this is handler func for task:", params[0].(int), params[1].(int), params[2].(int))
	// 	params[3].(chan int) <- params[0].(int) + params[1].(int) + params[2].(int)

	// 	// for i := 0; i < len(params); i++ {
	// 	// 	myFunc(int32(params[i].(int)))
	// 	// }
	// 	//time.Sleep(time.Second * 2)
	// 	wg.Done()
	// })
	// defer p.Release()
	// // Submit tasks one by one.
	// for i := 0; i < runTimes; i++ {
	// 	wg.Add(1)
	// 	_ = p.Invoke(i, i+4, i+8, resCh)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(resCh)
	// }()

	// sum := 0
	// for v := range resCh {
	// 	fmt.Println("result:", v)
	// 	sum += v
	// }
	// fmt.Printf("running goroutines: %d\n", p.Running())
	// fmt.Printf("finish all tasks, result is %d\n", sum)
	// fmt.Printf("sum 10: %d\n", sum)

	// 	// 3. Use the MultiPool and set the capacity of the 10 goroutine pools to unlimited.
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
