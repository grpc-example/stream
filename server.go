package main

import (
	"fmt"
	"github.com/grpc-example/stream/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

const PORT = ":50051"

type server struct {
	pb.UnimplementedGreeterServer
}

// GetStream 服务端 单向流
func (s *server) GetStream(req *pb.StreamReqData, res pb.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		err := res.Send(&pb.StreamResData{Data: fmt.Sprintf("%v", time.Now().Unix())})
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		if i > 10 {
			break
		}
	}
	return nil
}

// PutStream 客户端 单向流
func (s *server) PutStream(cliStr pb.Greeter_PutStreamServer) error {

	for {
		if tem, err := cliStr.Recv(); err == nil {
			log.Println(tem)
		} else {
			log.Println("break, err :", err)
			break
		}
	}

	return nil
}

// AllStream 客户端服务端 双向流
func (s *server) AllStream(allStr pb.Greeter_AllStreamServer) error {

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			data, _ := allStr.Recv()
			log.Println(data)
		}
		wg.Done()
	}()

	go func() {
		for {
			err := allStr.Send(&pb.StreamResData{Data: "ssss"})
			if err != nil {
				return
			}
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func main() {
	//监听端口
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		return
	}
	//创建一个grpc 服务器
	s := grpc.NewServer()
	//注册事件
	pb.RegisterGreeterServer(s, &server{})
	//注意这里这个pro是你定义proto文件生成后的那个go文件中引用的，而后面的这个函数是注册名称，是根据你自己定义的名称生成的，不同的文件该函数名称是不一样的，不过都是register这个单词开头的，这里不能原样照搬
	//处理链接
	log.Fatal(s.Serve(lis))
}
