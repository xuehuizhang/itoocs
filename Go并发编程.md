### Go语言并发管道与Go程



 #### 并发/并行概念

并发：同一个时间段内执行多个任务 (单核CPU执行多个任务)

并行：同一时刻同时执行多个任务（多核CPU执行多个任务）

Go语言实现并发通过Goroutine实现

Goroutine属于用户态线程，是由Go语言运行时调度完成，并不像java/c#等语言中用内核态线程实现并发，Goroutine的调度完全在用户态完成，这样就避免了java/c#中实现并发内核态与用户态之间的频繁切换

#### 案例1

```go
package main
import (
	"fmt"
	"sync"
)
var (
	wg sync.WaitGroup //这里通过WaitGrop等待Goroutine并发执行完成
)
func hello(i int) {
	defer wg.Done() //defer要先注册，不然，后续出了异常，defer就不会在执行了
	fmt.Println("hello=", i)
}
func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go hello(i)
	}
	wg.Wait()
}
```

#### 

#### Goroutine与线程

OS线程（操作系统线程）一般都有固定的栈内存（2M）

Goroutine的栈内存是可变的，2k-1GB

#### Goroutine调度

MPG调度模式

#### GOMAXPROCS

案例2

```go
package main
import (
	"fmt"
	"runtime"
	"time"
)
func printHello() {
	for i := 0; i < 1000; i++ {
		fmt.Println("hello--", i)
	}
}
func printWorld() {
	for i := 0; i < 100; i++ {
		fmt.Println("---------------world--", i)
	}
}
func main() {
	runtime.GOMAXPROCS(8)
	go printHello()
	go printWorld()
	time.Sleep(time.Second)
}
```

Go语言中的操作系统线程和goroutine的关系：

1. 一个操作系统线程对应用户态多个goroutine。
2. go程序可以同时使用多个操作系统线程。
3. goroutine和OS线程是多对多的关系，即m:n。

### Channel

在其他语言中都是通过<u>共享内存的方式来进行通信</u>，但是在Go语言中提倡通过<u>通信的方式来共享内存</u>

GO语言的并发模型CSP

channel是一种特殊的数据类型，类似于队列，先入先出

##### channel类型

是一种引用类型，零值为nil

声明：

	var 变量名 chan 元素类型

	var c1 chan int

初始化

	channel声明后要进行初始化，不然默认为nil

	c1=make(chan int) 无缓冲的

	c1=make(chan int,1)有缓冲

案例3

```go
package main

func main() {
	var c1 chan int
	var c2 chan string
	c1 = make(chan int)
	c2 = make(chan string, 1)
}
```



##### channel操作

发送  <-

```go
c1 <- 10
```



接收  <-

```go
n:=<-c1
```



关闭  close

```go
close(c1)
```

关闭后的通道有以下特点：

1. 对一个关闭的通道再发送值就会导致panic。
2. 对一个关闭的通道进行接收会一直获取值直到通道为空。
3. 对一个关闭的并且没有值的通道执行接收操作会得到对应类型的零值。
4. 关闭一个已经关闭的通道会导致panic。

#### 无缓冲通道

案例4

```go
package main

import "fmt"

func main() {
	c := make(chan int)
	c <- 10
	fmt.Println("写入成功")
}
```

案例4的代码在编译的时候是没问题的，但是在运行的时候会报错，原因是无缓冲的通道，只有存在接收的时候，才能写入数据，否则会死锁，一直等待接收就位才会写入数据，解决办法

案例5

```go
package main

import "fmt"

func recv(c1 chan int) {
	n := <-c1
	fmt.Println("接收成功：", n)
}
func main() {
	c := make(chan int)
	go recv(c)    //要放在写入之前
	c <- 10
	fmt.Println("写入成功")
}
```

#### 有缓冲通道

案例6

```go
package main

import "fmt"

func main() {
	c := make(chan int, 1)
	c <- 10
	fmt.Println("写入成功")
}
```

案例5采用的是有缓冲的通道，在不超过缓冲大小的情况下写入数据，就不会存在案例4中存在的问题

#### 关闭通道Close

案例6 

```go
package main

import (
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func sendData(c1 chan int) {
	fmt.Println("send data")
	for i := 0; i < 10; i++ {
		c1 <- i
	}
	close(c1)
	wg.Done()
}

func readData(c2 chan int, c1 chan int) {
	fmt.Println("read data")
	for {
		v, ok := <-c1
		if !ok {
			fmt.Println("读取数据完毕")
			break
		}
		c2 <- v
	}
	close(c2)
	wg.Done()
}

func main() {
	var c1 = make(chan int)
	var c2 = make(chan int)
	wg.Add(2)
	go sendData(c1)
	go readData(c2, c1)

	for v := range c2 {
		fmt.Println(v)
	}
	wg.Wait()
	fmt.Println("main")
}
```

案例6中使用了两种判断通道是否关闭 v, ok := <-c1    for-range(常用)

#### 单向通道

单向通道一般用户参数传递当作形参

双向通道可以转换为单向通道，单向通道不能转换成双向通道

案例7

```go
package main

import (
	"fmt"
	"sync"
)

var (
	wg sync.WaitGroup
)

func sendData(c1 chan<- int) {
	fmt.Println("send data")
	for i := 0; i < 10; i++ {
		c1 <- i
	}
	close(c1)
	wg.Done()
}

func readData(c2 chan<- int, c1 <-chan int) {
	fmt.Println("read data")
	for {
		v, ok := <-c1
		if !ok {
			fmt.Println("读取数据完毕")
			break
		}
		c2 <- v
	}
	close(c2)
	wg.Done()
}

func main() {
	var c1 = make(chan int)
	var c2 = make(chan int)
	wg.Add(2)
	go sendData(c1)
	go readData(c2, c1)

	for v := range c2 {
		fmt.Println(v)
	}
	wg.Wait()
	fmt.Println("main")
}
```

#### worker pool

作用：防止开启过多的go程

案例8

```go
package main

import (
	"fmt"
	"time"
)

func worker(i int, jobs chan int, result chan int) {
	for job := range jobs {
		fmt.Printf("worker %d 开始处理 job %d\n", i, job)
		time.Sleep(time.Second)
		fmt.Printf("worker %d 处理完毕 job %d\n", i, job)
		result <- job * job
	}
}

func main() {
	var jobs = make(chan int, 10)
	var result = make(chan int)

	for i := 0; i < 3; i++ {
		go worker(i, jobs, result)
	}

	for i := 0; i < 10; i++ {
		jobs <- i
	}
	close(jobs)
	for i := 0; i < 10; i++ {
		fmt.Println("result:", <-result)
	}
}

```

#### select

案例9

```go
package main

import "fmt"

func main() {
	var ch = make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case n := <-ch:
			fmt.Println(n)
		case ch <- i:
		}
	}
}
```

使用`select`语句能提高代码的可读性。

如果多个`case`同时满足，`select`会随机选择一个。

对于没有`case`的`select{}`会一直等待，可用于阻塞main函数。







