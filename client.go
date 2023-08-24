package main

import (
	"github.com/grpc-example/stream/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"context"
	_ "google.golang.org/grpc/balancer/grpclb"
	"log"
	"time"
)

const (
	ADDRESS = "localhost:50051"
)

func main() {
	//通过grpc 库 建立一个连接
	conn, err := grpc.Dial(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()
	//通过刚刚的连接 生成一个client对象。
	c := pb.NewGreeterClient(conn) //这个pro与服务端的同理，也是来源于proto编译生成的那个go文件内部调用
	//调用服务端推送流
	reqStreamData := &pb.StreamReqData{Data: "aaa"}
	res, _ := c.GetStream(context.Background(), reqStreamData)
	for {
		aa, err := res.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(aa)
	}
	//客户端 推送 流
	putRes, _ := c.PutStream(context.Background())
	i := 1
	for {
		i++
		putRes.Send(&pb.StreamReqData{Data: "ss"})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}
	//服务端 客户端 双向流
	allStr, _ := c.AllStream(context.Background())
	go func() {
		for {
			data, err := allStr.Recv()
			if err != nil {
				log.Println(err)
				break
			}
			log.Println(data)
		}
	}()

	go func() {
		for {
			allStr.Send(&pb.StreamReqData{Data: "ssss"})
			time.Sleep(time.Second)
		}
	}()
	select {}
}
