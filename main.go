/*
 * Author Linvon
 */

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	fileList := parseArgs()

	t := time.Now()
	initData()
	for _, v := range fileList {
		s := time.Now()
		file, isGz, err := getFile(v)
		if err != nil {
			fmt.Printf("cannot able to read the file %v, err %v\n", v, err)
			continue
		}
		err = process(file, isGz)
		if err != nil {
			fmt.Printf("handle filr %v err %v\n", v, err)
		}
		_ = file.Close()
		fmt.Printf("File %v Time taken - %v\n", v, time.Since(s))
	}
	fmt.Printf("All Time taken - %v\n", time.Since(t))
	outPut()
}

// TODO Customize your log file format : xlog-2020-07-27_00001(.gz)
func parseArgs() []string {
	args := os.Args[1:]
	if len(args) != 3 {
		log.Fatalln("Format ./coreader 200727 0 190(date start end)")
	}
	format := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalln("Format ./coreader 200727 0 190(date start end)")
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalln("Format ./coreader 200727 0 190(date start end)")
	}

	fList := make([]string, 0)
	for i := start; i < end; i++ {
		s := fmt.Sprintf(format, i)
		fList = append(fList, s)
	}
	return fList
}
