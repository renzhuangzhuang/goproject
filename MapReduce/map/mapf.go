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
	/* path := "../client/" + filename
	err = os.Remove(path)
	if err != nil {
		log.Fatal("删除失败: ", err)
	} */
}
