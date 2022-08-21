package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mapreduce/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := protos.NewMrServiceClient(conn)
	stream, err := client.GetSStream(context.Background(), &protos.MrRequest{Data: []byte{}})
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	for {
		count++
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
	}
}
