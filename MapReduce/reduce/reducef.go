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
