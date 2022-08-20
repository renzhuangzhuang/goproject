package main

import (
	"context"
	"fmt"
	"grpcdemo/protos"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
)

type demoService struct {
	protos.UnimplementedDemoServiceServer
}

// 实现方法
func (d *demoService) GetDemo(ctx context.Context, req *protos.DemoRequst) (*protos.DemoResponse, error) {
	res := req.GetName()
	fmt.Println("接收到的姓名为: ", res)
	return &protos.DemoResponse{Name: res}, nil

}

// 客户端流
func (d *demoService) GetCStream(stream protos.DemoService_GetCStreamServer) error {
	//
	count := 0 //用于计数表示客户端发送结束
	for {
		res, err := stream.Recv() // 一直接收数据
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Println("服务端接收到的数据流: ", res.Name, count)
		count++
		if count > 10 {
			req := protos.DemoResponse{Name: "接收完毕"}
			err := stream.SendAndClose(&req) // 发送信息给客户端
			if err != nil {
				return err
			}
			return nil

		}
	}

}

// 服务端流
func (d *demoService) GetSStream(requst *protos.DemoRequst, stream protos.DemoService_GetSStreamServer) error {
	count := 0
	for {
		rsp := protos.DemoResponse{Name: requst.Name} // 那个一个请求标志
		err := stream.Send(&rsp)
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
		count++
		if count > 10 {
			return nil
		}

	}
}

// 双端流
func (d *demoService) GetBiStream(stream protos.DemoService_GetBiStreamServer) error {
	for {
		res, err := stream.Recv()
		if err != nil {
			return nil
		}
		fmt.Println("服务端收到客户端消息", res.Name)
		time.Sleep(time.Second)
		rsp := &protos.DemoResponse{Name: "向客户端发送数据"}
		err = stream.Send(rsp)
		if err != nil {
			return nil
		}
	}
}

func main() {
	rpcs := grpc.NewServer()
	protos.RegisterDemoServiceServer(rpcs, &demoService{})
	listen, err := net.Listen("tcp", ":8080")

	if err != nil {
		panic(err)
	}
	rpcs.Serve(listen)
	defer listen.Close()
}
