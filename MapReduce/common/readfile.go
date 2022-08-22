package readfile

import (
	"log"
	"math"
	"os"
	"strconv"
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

func ReduceName(mapTask int) string {
	return "mrclient." + "-" + strconv.Itoa(mapTask) + ".txt"
}
func ReduceName_Json(mapTask int) string {
	return "mrclient." + "-" + strconv.Itoa(mapTask) + ".json"
}

func Task_n(name string, chunksize int) int {
	fi, err := os.Stat(name)
	if err != nil {
		log.Fatal("获取文件大小失败,", err)
		return -1
	}
	file_num := math.Ceil(float64(fi.Size()) / float64(chunksize))
	return int(file_num)

}

func RemoveFile(task int, reduce_name func(mapTask int) string) {
	file_name := reduce_name(task)
	path := "../client/" + file_name
	err := os.Remove(path)
	if err != nil {
		log.Fatal("删除失败: ", err)
	}
}
