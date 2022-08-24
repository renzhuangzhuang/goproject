## 基于GPRC实现MapReduce

### 1、实现文件切片划分

***文件切分***

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

### 2、 配置GRPC

使用grpc框架实现多机之间网络传输

#### 2.1第一步编辑proto文件

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

这边定义两种数据类型，分别是发送数据和接收数据二者都是比特类型，服务定义了两种`GetSStream`服务端流、`GetBiStream`双端流。

#### 2.2第二步定义服务端和客户端

```go
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
```

在第一步中，采用指令：

```sh
protoc --go_out=.(文件相对路径) --go_opt=paths=source_relative --go-grpc_out=.(文件相对路径) --go-grpc_opt=paths=source_relative 文件名.proto
```

将会得到两个go文件分别为：.._grpc.pb.go和.pb.go文件，里面定义了相关服务的方法，因此在服务端我们只要实例化结构体，并实现接口中的方法即可使用对应的API，同样还包含定义的信息类型结构体，在proto文件中定义了一个`bytes`类型的变量，在具体的go结构体中就是一个比特切片，

在服务端实现了两种方法--对应服务端流和双端流，两种模式如下图所示：

![服务端流](D:\goproject\src\MapReduce\服务端流.png)

![双端流模式](D:\goproject\src\MapReduce\双端流模式.png)

***服务端流***：根据图片内容可知，我在实现服务端流时候，客户端首先会发送一个空数据包过来，触发服务端请求，然后再服务端对文件进行切片处理并分块发送，在客户端处将得到数据保存在本地磁盘，这步操作是为了避免因特殊原因导致客户端宕机导致内存中数据丢失。

***双端流***：和服务端流设计思路类似，在建立连接后，客户端首先发送一条空数据过来（当然可以发送任意类型的数据）在服务端接收数据后，对数据进行判断，若是空数据，则表示服务端需要向客户端发送处理好的切片文件，再发送完切片数据后，需要想客户端发送空包提醒客户端不需要再接收服务端的信息（因为这时候服务端发送任务已经结束，而客户端采用for循环接收数据，若没有数据则会一直原地等待），剩余操作即为客户端处理发送来的文件，最后将处理好的数据发送给服务端。

客户端代码实现：

```go
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	readfile "mapreduce/common"
	mapf "mapreduce/map"
	"mapreduce/protos"
	reducef "mapreduce/reduce"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var wg sync.WaitGroup
var wg1 sync.WaitGroup

const chunk = 1 << 12

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := protos.NewMrServiceClient(conn)
	stream, err := client.GetBiStream(context.Background())
	if err != nil {
		log.Fatal("连接失败: ", err)
	}
	mapTask := 0
	requst := &protos.MrRequest{Data: []byte{}}
	err = stream.Send(requst)
	if err != nil {
		log.Fatal("发送失败: ", err)
	}
	for {
		fmt.Println("等待接收数据")
		res, err := stream.Recv()
		if err != nil {
			log.Fatal("出错,", err)
		}
		if res.Data == nil {
			break
		}
		filename := readfile.ReduceName(mapTask)
		ioutil.WriteFile(filename, res.Data, 0666)
		mapTask++
	}
	time.Sleep(time.Second * 10)

	//
	fmt.Println("处理数据: ")
	JSON_num := 0
	for i := 0; i < mapTask; i++ {
		filename := readfile.ReduceName(i)
		ans := make(map[string]int)
		ReduceTask := readfile.Task_n(filename, chunk)
		wg.Add(1)
		go mapf.DoMap(chunk, filename, int(ReduceTask), ans, &wg)
		result := make([]chan map[string]int, int(ReduceTask))
		JSON_num = ReduceTask
		for i := 0; i < int(ReduceTask); i++ {
			result[i] = make(chan map[string]int, 10000)
			wg1.Add(1)
			file_name := readfile.ReduceName_Json(i)
			go reducef.DoReduce(file_name, result[i], &wg1)
		}
		wg.Wait()
		wg1.Wait()
		result_all := make(map[string]int)
		for _, v := range result {
			for k, value := range <-v {
				result_all[k] += value
			}
			if len(v) == 0 {
				close(v)
			}
		}
		sortmap := []string{}
		for k := range result_all {
			sortmap = append(sortmap, k)
		}
		sort.Strings(sortmap)
		//保存结果
		final_result, err := os.Create("result" + strconv.Itoa(i) + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		defer final_result.Close()

		for _, v := range sortmap {
			final_result.WriteString(v + ":" + strconv.Itoa(result_all[v]) + "\n")
		}

	}
	//给服务端发送处理后的结果
	for i := 0; i < mapTask; i++ {
		fname := readfile.ReduceName_res(i)
		file, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		stats, err := file.Stat()
		if err != nil {
			log.Fatal("获取文件大小失败: ", err)
		}
		data := make([]byte, stats.Size())
		_, err = file.Read(data)
		if err != nil {
			log.Fatal(err)
		}
		requst := &protos.MrRequest{Data: data}
		err = stream.Send(requst)
		if err != nil {
			log.Fatal(err)
		}
	}
	requst = &protos.MrRequest{Data: []byte("close")}
	stream.Send(requst)

	fmt.Println("删除多余文件")
	for i := 0; i < mapTask; i++ {
		readfile.RemoveFile(i, readfile.ReduceName)
		//readfile.RemoveFile(i, readfile.ReduceName_res)
	}
	for j := 0; j < JSON_num; j++ {
		readfile.RemoveFile(j, readfile.ReduceName_Json)
	}
	/* for i := 0; i < mapTask; i++ {
		readfile.RemoveFile(i, readfile.ReduceName_res)
	} */

}

```

### 3、设计细节和思路

整体思路：

![工作流程](D:\goproject\src\MapReduce\工作流程.png)

**建立连接**

​		建立连接即通过GRPC框架，服务端持续监听主机的某个端口，客户端访问相应端口建立连接，代码如下:

```go 
//服务端
rpcs := grpc.NewServer()
protos.RegisterMrServiceServer(rpcs, &MrService{})
listen, err := net.Listen("tcp", ":8080")
if err != nil {
    panic(err)
}
pcs.Serve(listen)
defer listen.Close()
//客户端
conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
```

​       因为程序模拟的是在本地运行所以只要在IP地址处填写本地IP即可。

**实现双端流方法**

         ```go
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
         ```

这里需要强调一点就是，发送和接收数据可以进行多次，比方说跳出for循环依然可以发送接收。在for循环中并没有终止条件，这里这样设计我考虑就是通过客户端来决定和服务端何时断开连接，而服务端在开启时就一直保持监听状态。

**客户端里的MapReduce操作**

直接上代码

```go
package mapf

import (
	"encoding/json"
	"fmt"
	"log"
	readfile "mapreduce/common"
	"os"
	"strings"
	"sync"
)

// 对得到的文件继续划分
func DoMap(
	chunksize int, //缓存大小
	filename string, // 文件名称
	nReduceTask int, //reduce任务数，块数
	ans map[string]int, // 保存中间存放k-v
	wg *sync.WaitGroup,
) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("读取文件出错:", err)
	}
	defer file.Close()
	//b1 := bufio.NewReader(file)

	for i := 0; i < nReduceTask; i++ {
		size := i * chunksize
		p := make([]byte, chunksize)
		n1, err := file.ReadAt(p, int64(size))
		if err != nil {
			log.Fatal("文件分块出错:", err)
		}
		ss := strings.Fields(string(p[:n1]))
		for _, v := range ss {
			word := strings.ToLower(v)
			for len(word) > 0 && (word[0] < 'a' || word[0] > 'z') {
				word = word[1:]
			}
			for len(word) > 0 && (word[len(word)-1] < 'a' || word[len(word)-1] > 'z') {
				word = word[:len(word)-1]
			}
			ans[word]++
		}
		file_name := readfile.ReduceName_Json(i)
		files, err := os.Create(file_name)
		if err != nil {
			log.Fatal("文件生成失败:", err)
			os.Exit(1)
		}
		taskJson, err := json.Marshal(ans)
		if err != nil {
			fmt.Println("错误一")
			panic(err)
		}
		if _, err := files.Write(taskJson); err != nil {
			fmt.Println("错误二")
			panic(err)
		}
		if err := files.Close(); err != nil {
			fmt.Println("错误三")
			panic(err)
		}
	}
	defer wg.Done()
}
```

整体思路即为：首先读取本地文件，并对文件进行切片操作，对分块好数据进行筛选（剔除特殊符号）,之后生成JSON文件（同样可以生成其他形式）。

```go
package reducef

import (
	"encoding/json"
	"os"
	"sync"
)

func DoReduce(
	filename string,
	result chan map[string]int,
	wg *sync.WaitGroup,
) {
	files, _ := os.Open(filename)
	defer files.Close()
	ans := make(map[string]int)
	ans1 := make(map[string]int)
	dec := json.NewDecoder(files)
	dec.Decode(&ans)
	for k, v := range ans {
		ans1[k] += v
	}
	result <- ans1
	defer wg.Done()
}
```

reduce操作则是构建一个map缓存，key为string类型，value为int类型。

设计map 和reduce 是可以支持并发操作的。

