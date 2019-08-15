# itoocs
### 类库项目
    1, go协程池 workerPool
       tag: 
            jobsChan:    存储任务的chan
            resultChan:  存储结果的chan
            workerPool 在jobsChan读取到任务后用于分配工人，也就是把任务分配协程去处理
            worker用于真正的去处理业务，处理结果放到resultChan中
            
### GO并发系列
    1，Go并发之管道与Go程
