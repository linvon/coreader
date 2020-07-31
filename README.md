# coreader
High efficiency big log file reader/高效率大日志文件处理工具

[中文文档](#是什么)

## What is this?

A tool to help you analysis big log file

## Feature
- Support normal file and gz file
- Replace std json with [json-iterator](https://github.com/json-iterator/go), improve unmarshall efficiency
- Use **sync.Pool** to reduce the pressure on GC
- Concurrently process the buffer chunk
- Read GBs file by lines in seconds, read in tens of seconds if with unmarshall

## How to use

`go clone https://github.com/linvon/coreader.git`

##### Customize your format function

Usually log file was named `xlog-2020-07-27_00001`, and compressed to `xlog-2020-07-27_00001.gz`

So input fotmat like`xlog-2020-07-27_00%03d` and start, end like 0 100, then it will handle file from `xlog-2020-07-27_00001` to `xlog-2020-07-27_00100`, also you can handle .gz file

You can customize your own format function if needed, return a filename list

``` go
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
```

##### Customize your record data and process function

In customize.go you can customize how to process your log file

``` go
var gMu sync.Mutex			
var allMap, errMap map[string]int64		// Your stat map, for easily output

// TODO Init your global record data
func initData() {
  allMap = make(map[string]int64, 0)
	errMap = make(map[string]int64, 0)
}

// TODO complete your line process function
func lineProcess(line string) {
	// handle log line here

	// Update Global Data, remember to lock
	gMu.Lock()
	allMap[log.HDid] ++
	if log.HAv != "" {
		errMap[log.HDid] ++
	}
	gMu.Unlock()

}

// TODO Customize your Output format
func outPut() {
  // output format you want
	s, _ := json.Marshal(allMap)
	fmt.Println(string(s))
	s, _ = json.Marshal(errMap)
	fmt.Println(string(s))
}

// If your log is in JSON format, you can process it like below

// TODO complete process function of unmarshal if needed
func UnmarshalLog(line string) (log LogSt, err error) {
	// pre process

	err = json.Unmarshal([]byte(line), &log)
	if err != nil {
		return
	}

	// post process

	return
}

// TODO your LogSt if needed
type LogSt struct {
	HAv    string `json:"h_av"`
	HDt    int    `json:"h_dt"`
	HOs    int    `json:"h_os"`
	HApp   string `json:"h_app"`
	HModel string `json:"h_model"`
	HDid   string `json:"h_did"`
	HNt    int    `json:"h_nt"`
	HCh    string `json:"h_ch"`
	HTs    int64  `json:"h_ts"`
	HLang  string `json:"h_lang"`
	HPkg   string `json:"h_pkg"`
	HM     int    `json:"h_m"`
}

```

##### In the end

`cd coreader; go build`

`./coreader format start end`



## 是什么?

是一个帮助你处理大日志文件的工具

## 特性

- 支持常规文件和 gzip 压缩文件
- 使用 [json-iterator](https://github.com/json-iterator/go) 代替标准库的 json， 提升 unmarshall 效率
- 使用 **sync.Pool** 减小GC 压力
- 并发处理文件缓冲块
- 在秒级内完成GB级文件按行读取，附带 unmarshall 动作时在十秒级完成读取

## 如何使用

`go clone https://github.com/linvon/coreader.git`

##### 定制你的格式化函数

通常情况下日志文件都是命名类似于 `xlog-2020-07-27_00001`, 并且压缩为 `xlog-2020-07-27_00001.gz`

因此给程序输入类似于`xlog-2020-07-27_00%03d` 这样的 format，以及类似于 0 100 这样的 start 和 end ，之后程序就会处理从 `xlog-2020-07-27_00001` 到 `xlog-2020-07-27_00100`这些日志文件, 同样你也可以处理 .gz 为后缀的压缩文件

如果有其他需要，你可以自己定制格式化函数，返回一个文件名称列表即可

``` go
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
```

##### 定制你需要记录的数据和处理函数

在 customize.go 文件中你可以定制如何去处理日志文件

``` go
var gMu sync.Mutex			
var allMap, errMap map[string]int64		// 使用 map 记录数据，方便输出

// TODO Init your global record data
func initData() {
  allMap = make(map[string]int64, 0)
	errMap = make(map[string]int64, 0)
}

// TODO complete your line process function
func lineProcess(line string) {
	// 在这里处理每一行日志

	// 处理数据后，更新统计数据，记得要加锁
	gMu.Lock()
	allMap[log.HDid] ++
	if log.HAv != "" {
		errMap[log.HDid] ++
	}
	gMu.Unlock()

}

// TODO Customize your Output format
func outPut() {
  // 设置你需要的输出格式
	s, _ := json.Marshal(allMap)
	fmt.Println(string(s))
	s, _ = json.Marshal(errMap)
	fmt.Println(string(s))
}

// 如果你的日志是 JSON 格式的，你可以按照下面的方法处理

// TODO complete process function of unmarshal if needed
func UnmarshalLog(line string) (log LogSt, err error) {
	// pre process

	err = json.Unmarshal([]byte(line), &log)
	if err != nil {
		return
	}

	// post process

	return
}

// TODO your LogSt if needed
type LogSt struct {
	HAv    string `json:"h_av"`
	HDt    int    `json:"h_dt"`
	HOs    int    `json:"h_os"`
	HApp   string `json:"h_app"`
	HModel string `json:"h_model"`
	HDid   string `json:"h_did"`
	HNt    int    `json:"h_nt"`
	HCh    string `json:"h_ch"`
	HTs    int64  `json:"h_ts"`
	HLang  string `json:"h_lang"`
	HPkg   string `json:"h_pkg"`
	HM     int    `json:"h_m"`
}

```

##### 最后

`cd coreader; go build`

`./coreader format start end`


