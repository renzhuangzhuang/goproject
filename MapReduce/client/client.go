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
