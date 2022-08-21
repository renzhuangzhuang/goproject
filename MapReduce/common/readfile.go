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
