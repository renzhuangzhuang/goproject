## 基于GPRC实现MapReduce

### 1、实现文件切片划分、多机之间传输、本地磁盘的保存

#### 1.1文件切分

关于文件切片划分在我的github上介绍了几种方式，[](https://github.com/renzhuangzhuang/goNotes/tree/master)，这里我就使用的`bufio.NewReader`+计算分块大小方式

在common文件夹下实现文件切片功能：

```go
package readfile

import (
	"math"
	"os"
)

//const filename string = "./file.txt"
//const chunksize = 1 << 10 // 块数可以自由划分

func Read_file(file string, chunksize int) ([][]byte, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	file_size := math.Ceil(float64(fi.Size()) / float64(chunksize)) // 得到文件块数
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buff := make([][]byte, int(file_size))
	for i := 0; i < int(file_size); i++ {
		size := i * chunksize
		buff[i] = make([]byte, chunksize)
		_, _ = f.ReadAt(buff[i], int64(size))
	}
	return buff, nil
}

```

生成适合网络传输的byte切片

#### 1.2 配置网络

使用grpc框架实现多机之间网络传输

第一步编辑proto文件

```protobuf
syntax = "proto3";

option go_package = "mapreduce/protos";

//定义发送消息
message MrRequest {
    bytes data = 1;
}

//定义接收消息
message MrResponse {
    bytes data = 1;
}

service MrService {
    //定义双端流
    rpc GetSStream(MrRequest) returns (stream MrResponse);
    rpc GetBiStream(stream MrRequest) returns (stream MrResponse);
}

```

第二步定义服务端和客户端

```go
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

```

```go
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

```

先开启服务端，再开启客户端，测试结果

```txt
g roll of parchment.\r\n\xa1\xa1\xa1\xa1\"When I call your name, you will put on the hat and sit on the stool to be sorted,\" she said. \"Abbott, Hannah!\"\r\n\xa1\xa1\xa1\xa1A pink-faced girl with blonde pigtails stumbled out of line, put on the hat, which fell right down over her eyes, and sat down. A moments pause --\r\n\xa1\xa1\xa1\xa1\"HUFFLEPUFF!\" shouted the hat.\r\n\xa1\xa1\xa1\xa1The table on the right cheered and clapped as Hannah went to sit down at the Hufflepuff table. Harry saw the ghost of the Fat Friar waving merrily at her.\r\n\xa1\xa1\xa1\xa1\"Bones, Susan!\"\r\n\xa1\xa1\xa1\xa1\"HUFFLEPUFF!\" shouted the hat again, and Susan scuttled off to sit next to Hannah.\r\n\xa1\xa1\xa1\xa1\"Boot, Terry!\"\r\n\xa1\xa1\xa1\xa1\"RAVENCLAW!\"\r\n\xa1\xa1\xa1\xa1The table second from the left clapped this time; several Ravenclaws stood up to shake hands with Terry as he joined them.\r\n\xa1\xa1\xa1\xa1\" Brocklehurst, Mandy\" went to Ravenclaw too, but \"Brown, Lavender\" became the first new Gryffindor, and the table on the far left exploded with cheers; Harry could see Ron's twin brothers catcalling.\r\n\xa1\xa1\xa1\xa1\"Bulstrode....
```

可以将文件进行切块传输