### Go并发安全与锁

#### 引入

案例1

```go
package main

import (
	"fmt"
	"sync"
)

var (
	x  int
	wg sync.WaitGroup
)

func add() {
	for i := 0; i <= 5000; i++ {
		x++
	}
	wg.Done()
}
func main() {
	wg.Add(2)
	go add()
	go add()
	wg.Wait()
	fmt.Println(x)
}
```

![](E:\GoIT\note\img\微信截图_20190815182825.png)

执行结果如上，由于并发共享资源导致计算结果错误

#### 互斥锁

sync.Mutex 要求同时只有一个线程去访问共享资源

案例2

```go
var (
	x    int
	wg   sync.WaitGroup
	lock sync.Mutex
)

func add() {
	for i := 0; i < 5000; i++ {
		lock.Lock()
		x++
		lock.Unlock()
	}
	wg.Done()
}
```

使用互斥锁能够保证同一时间有且只有一个`goroutine`进入临界区，其他的`goroutine`则在等待锁；当互斥锁释放后，等待的`goroutine`才可以获取锁进入临界区，多个`goroutine`同时等待一个锁时，唤醒的策略是随机的。

#### 读写互斥锁

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg     sync.WaitGroup
	lock   sync.Mutex
	rwlock sync.RWMutex
	n      int
)

func write() {
	//rwlock.Lock()
	lock.Lock()
	n++
	time.Sleep(time.Millisecond * 10)
	lock.Unlock()
	//rwlock.Unlock()
	wg.Done()
}

func read() {
	//rwlock.RLock()
	lock.Lock()
	time.Sleep(time.Millisecond)
	lock.Unlock()
	//rwlock.RUnlock()
	wg.Done()
}
func main() {
	start := time.Now()
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go write()
	}

	for i := 1; i < 1000; i++ {
		wg.Add(1)
		go read()
	}
	wg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(start))
}
```

![](E:\GoIT\note\img\1565881238.png)

案例2测试得出，在读多写少的情况下，读写锁是非常适合的

#### sync.WaitGroup

Go语言使用sync.WaitGroup实现并发同步问题，毕竟使用time.Sleep()的方式不合适

注意sync.WaitGroup是一个结构体，在当作参数传递的时候，必须传递指针

#### sync.Once

```go
package main

import (
	"fmt"
	"sync"
)

var data map[int]string
var wg sync.WaitGroup
var once sync.Once

func loadData() {
	fmt.Println("load data once")
	data = map[int]string{1: "data1", 2: "data2", 3: "data3"}
}

func initData() {
	// if data == nil {
	// 	loadData()       //使用这种方式 initData并发不安全的，导致loadData执行多次
	// }
	once.Do(loadData)    //这种方式load只执行一次
	fmt.Println("data:", data[1])
	wg.Done()
}

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go initData()
	}
	wg.Wait()
}
```

sync.Once 实现了单例模式

#### sync.Map

go语言中内置的map不是并发安全的

案例3

```go
package main

import (
	"fmt"
	"strconv"
	"sync"
)

var m = make(map[string]int, 10)  //并发不安全的
var wg sync.WaitGroup

func get(key string) int {
	return m[key]
}

func set(key string, n int) {
	m[key] = n
}

func main() {
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(a int) {
			key := strconv.Itoa(a)
			set(key, a+1)
			fmt.Printf("value is %d", get(key))
            wg.Done()
		}(i)
	}
	wg.Wait()
}

```

几个`goroutine`的时候可能没什么问题，当并发多了之后执行上面的代码就会报`fatal error: concurrent map writes`错误。

像这种并发场景下就可以考虑使用sync包中提供的map了

案例4

```go
package main

import (
	"fmt"
	"strconv"
	"sync"
)

var m sync.Map
var wg sync.WaitGroup

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			key := strconv.Itoa(n)
			m.Store(key, n*n)
			value, _ := m.Load(key)
			fmt.Println("get value :", value)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```

sync.Map是一个开箱即用的，并发安全的，内置了Store,Load,LoadOrStore,Delete,Range等方法

#### 原子操作

由于加锁操作涉及到内核态的上下文切换，比较耗时，因此可以考虑使用原子操作来保证并发安全，因为原子操作是Go语言提供的方法他在用户态就可以完成，因此性能比较好，由内置的标准库sync/atomic提供

案例5

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	x    int64
	lock sync.Mutex
	wg   sync.WaitGroup
)

func add() {
	x++
	wg.Done()
}

func lockAdd() {
	lock.Lock()
	x++
	lock.Unlock()
	wg.Done()
}

func atomicAdd() {
	atomic.AddInt64(&x, 1)
	wg.Done()

}

func main() {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go atomicAdd()
	}
	wg.Wait()
	end := time.Now()
	fmt.Println("cost", end.Sub(start))
	fmt.Println(x)
}
```

