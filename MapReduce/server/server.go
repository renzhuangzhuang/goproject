package main

import (
	"fmt"
	"io/ioutil"
	"log"
	readfile "mapreduce/common"
	"mapreduce/protos"
	"net"

	"google.golang.org/grpc"
)

const chunksize = 1 << 14

type MrService struct {
	protos.UnimplementedMrServiceServer
}

const filename = "../file.txt"

func (m *MrService) GetSStream(req *protos.MrRequest, stream protos.MrService_GetSStreamServer) error {
	buff, err := readfile.Read_file(filename, chunksize)
	if err != nil {
		log.Fatal("读取文件失败: ", err)
	}
	n := len(buff)
	for i := 0; i < n; i++ {
		rsp := protos.MrResponse{Data: buff[i]}
		err := stream.Send(&rsp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MrService) GetBiStream(stream protos.MrService_GetBiStreamServer) error {
	buff, err := readfile.Read_file(filename, chunksize)
	if err != nil {
		log.Fatal("读取文件失败: ", err)
	}
	n := len(buff)
	count := 0
	for {
		res, err := stream.Recv()
		if err != nil {
			return nil
		}
		//持续接收
		if res.Data != nil && string(res.Data) != "close" {
			count++
			fmt.Println("服务端收到客户端消息") // 这个就是一个result值
			filename := readfile.ReduceName(count)
			ioutil.WriteFile(filename, res.Data, 0666)
		} else {
			//空包 接收服务端数据
			for i := 0; i < n; i++ {
				rsp := &protos.MrResponse{Data: buff[i]}
				stream.Send(rsp)
			}
			rsp := &protos.MrResponse{Data: []byte{}}
			stream.Send(rsp)
		}
		if string(res.Data) == "close" {
			fmt.Println("合并数据")
		}
	}
}

func main() {
	rpcs := grpc.NewServer()
	protos.RegisterMrServiceServer(rpcs, &MrService{})
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	rpcs.Serve(listen)
	defer listen.Close()
}
