### RPC与GRPC

#### 原生rpc

什么是rpc：即远程过程调用，像调用本地函数那么简单调用远程函数

go内置的rpc包实现

案例1

服务端

```go
package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type Panda int

//GetInfo 获取详情
func (p *Panda) GetInfo(argType int, replyType *int) error {
	fmt.Println("打印客户端发过来的内容为：", argType)
	*replyType = argType + 1000
	return nil
}

func main() {
	p := new(Panda)
	rpc.Register(p) //注册一个对象

	rpc.HandleHTTP() //将rpc连接到网络中

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	http.Serve(ln, nil)
}
```

客户端：

```go
package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	var replay int
	err = client.Call("Panda.GetInfo", 100, &replay) //调用rpc中的方法

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(replay)
}
```



#### grpc

```bash
go get google.golang.org/grpc
```

```protobuf
syntax = "proto3";

package myproto;

//定义服务

service HelloServer{
    rpc SayHello (HelloReq) returns (HelloResp){}
    rpc SayName(NameReq)returns(NameResp){}
}

message HelloReq{
    string Name=1;
}

message HelloResp{
    string Msg=1;
}

message NameReq{
    string Name=1;
}

message NameResp{
    string Msg=1;
}
```

服务端

```go
package main

import (
	"google.golang.org/grpc"
	"net"
	"fmt"
	pb "day03/myproto"
	"context"
)

//实现 HelloServer 中的方法
type server struct{
}


func(s *server)SayHello(ctx context.Context, in *pb.HelloReq) (out *pb.HelloResp,err error){
	return &pb.HelloResp{Msg:"hello"+in.Name},nil
}

func(s *server) SayName(ctx context.Context, in *pb.NameReq) (out *pb.NameResp,err error){
	return &pb.NameResp{Msg:in.Name+"你好"},nil
}

func main(){
	ln,err:=net.Listen("tcp",":10086")
	if err!=nil{
		fmt.Println("error:",err)
	}

	//创建grpc服务
	srv:=grpc.NewServer()
	//注册服务
	pb.RegisterHelloServerServer(srv,&server{}) 

	err=srv.Serve(ln)

	if err!=nil{
		fmt.Println("error:",err)
	}

}
```

客户端

```go
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pd "day03/myproto"
)


func main(){
	conn,err:=grpc.Dial("127.0.0.1:10086",grpc.WithInsecure())
	if err!=nil{
		fmt.Println("链接远程服务算错误")
		panic("error")
	}

	defer conn.Close()

	c:=pd.NewHelloServerClient(conn)

	re,err:=c.SayHello(context.Background(),&pd.HelloReq{Name:"考拉"})

	if err!=nil{
		fmt.Println("HelloReq Error")
		panic("error")
	}

	fmt.Println("Say Hello:",re.Msg)


	res,err:=c.SayName(context.Background(),&pd.NameReq{Name:"panda"})
	if err!=nil{
		fmt.Println("Name req error")
		panic(err)
	}
	fmt.Println("Name resp:",res.Msg)
}
```



