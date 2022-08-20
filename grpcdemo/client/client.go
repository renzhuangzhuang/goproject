package main

import (
	"context"
	"fmt"
	"grpcdemo/protos"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 客户端流 发送数据
func demoRequst(stream protos.DemoService_GetCStreamClient, rsp chan struct{}) {
	count := 0
	for {
		requst := protos.DemoRequst{
			Name: "count",
		}
		err := stream.Send(&requst)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		count++
		if count > 10 {
			rsp <- struct{}{}
			break
		}

	}
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials())) //  grpc.WithInsecure() 表示无认证 是不安全的 可以修改
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := protos.NewDemoServiceClient(conn)
	//response := &protos.DemoRequst{Name: "1"}
	stream, err := client.GetBiStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for {
		requst := &protos.DemoRequst{Name: "发送数据到服务器"}
		err = stream.Send(requst)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		res, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("客户端收到的信息: ", res.Name)
	}
	/* if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端数据接收完毕")
				err := stream.CloseSend()
				if err != nil {
					log.Fatal(err)
				}
				break
			}
			log.Fatal(err)
		}
		fmt.Println("客户端接收的流: ", res)
	} */
	/* if err != nil {
		log.Fatal(err)
	}
	rsp := make(chan struct{}, 1)
	go demoRequst(stream, rsp)

	select {
	case <-rsp:
		res, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatal(err)
		}
		res_message := res.Name
		fmt.Println(res_message)
	} */

	/* res, err := client.GetDemo(context.Background(), &protos.DemoRequst{Name: "任壮壮"})
	if err != nil {
		log.Fatal("调用方法出错: ", err)
	}
	fmt.Println("调用成功 ProdStock", res) */

}
