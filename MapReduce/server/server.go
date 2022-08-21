package main

import (
	"log"
	readfile "mapreduce/common"
	"mapreduce/protos"
	"net"

	"google.golang.org/grpc"
)

const chunksize = 1 << 20

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
