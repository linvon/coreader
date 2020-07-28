/*
 * Author Linvon
 */

package main

import (
	"fmt"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var gMu sync.Mutex
var allMap, errMap map[string]int64

// TODO Init your global record data
func initData() {
	allMap = make(map[string]int64, 0)
	errMap = make(map[string]int64, 0)
}

// TODO complete your line process function
func lineProcess(line string) {
	// pre process

	log, err := UnmarshalLog(line)
	if err != nil {
		return
	}

	// post process

	// Update Global Data 
	gMu.Lock()
	allMap[log.HDid] ++
	if log.HAv != "" {
		errMap[log.HDid] ++
	}
	gMu.Unlock()

}

// TODO Customize your Output format
func outPut() {
	s, _ := json.Marshal(allMap)
	fmt.Println(string(s))
	s, _ = json.Marshal(errMap)
	fmt.Println(string(s))
}

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
