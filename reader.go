/*
 * Author Linvon
 */

package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
)

func getFile(fileName string) (f *os.File, isGz bool, err error) {
	f, err = os.Open(fileName)
	if strings.HasSuffix(fileName, ".gz") {
		isGz = true
	}
	return
}

func process(f *os.File, isGz bool) error {
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	var r *bufio.Reader
	if isGz {
		fz, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		r = bufio.NewReader(fz) //解压成功后读取解压后的文件
	} else {
		r = bufio.NewReader(f) //解压失败（还是读取原来文件）gz文件还是读取原始文件
	}

	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]

		if n == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				break
			}
			return err
		}

		nextUntillNewline, err := r.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}

		wg.Add(1)
		go func() {
			processChunk(buf, &linesPool, &stringPool)
			wg.Done()
		}()

	}

	wg.Wait()
	return nil
}

func processChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool) {
	var wg2 sync.WaitGroup

	logs := stringPool.Get().(string)
	logs = string(chunk)

	linesPool.Put(chunk)

	logsSlice := strings.Split(logs, "\n")

	stringPool.Put(logs)

	chunkSize := 300
	n := len(logsSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}

	for i := 0; i < (noOfThread); i++ {
		wg2.Add(1)
		go func(s int, e int) {
			defer wg2.Done() //to avaoid deadlocks
			for i := s; i < e; i++ {
				text := logsSlice[i]
				if len(text) == 0 {
					continue
				}
				lineProcess(text)
			}
		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	}

	wg2.Wait()
	logsSlice = nil
}
