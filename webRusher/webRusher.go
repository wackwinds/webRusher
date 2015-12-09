// webRusher project webRusher.go
package webRusher

import (
	"net/url"
	"fmt"
	"log"
	"runtime"
	"sync"
	"strconv"
	"configReader"
	"fileReader"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
	"webAttacker"
)

var isDebug bool
var ok bool
var workerMutex sync.Mutex
var workerNum int64
var mainProcessChan chan bool
var threadContralChan chan bool

type webRusher struct {
	targetUrl string
	replayFilePath string
	startMark string
	endMark string
}

func init(){
	isDebug = false
	workerNum = 0;
	mainProcessChan = make(chan bool)
}

func (wRusher *webRusher) getInitConf(){
	conf := configReader.NewConfigReader("conf.txt")
	wRusher.targetUrl, ok = conf.GetConfig("targetUrl")
	if (!ok){
		log.Panic("failed to get target url")
	}
	
	wRusher.startMark, ok = conf.GetConfig("startMark")
	if (!ok){
		log.Panic("failed to get startMark")
	}
	
	wRusher.endMark, ok = conf.GetConfig("endMark")
	if (!ok){
		log.Panic("failed to get endMark")
	}
	
	wRusher.replayFilePath, ok = conf.GetConfig("replayFilePath")
	if (!ok){
		log.Panic("failed to get replay file path")
	}
	
	strCpuNumForThread, ok := conf.GetConfig("cpuNumForThread")
	
	// todo: deal with err 
	cpuNumForThread, _ := strconv.Atoi(strCpuNumForThread)
	
	if (!ok || (cpuNumForThread < 0 && cpuNumForThread != -1)){
		log.Println("failed to get cpuForThreadNum")
		cpuNumForThread = 0
	}
	
	switch cpuNumForThread{
		case -1:
		runtime.GOMAXPROCS(runtime.NumCPU())
		case 0:
		;
		default:
		runtime.GOMAXPROCS(cpuNumForThread)
	}
	
	strNumOfThread, ok := conf.GetConfig("threadNum")
	
	// todo: deal with err 
	numOfThread, _ := strconv.Atoi(strNumOfThread)
	
	if (!ok || (numOfThread <= 0)){
		log.Println("failed to get threadNum")
		threadContralChan = nil
	}else{
		threadContralChan = make(chan bool, numOfThread)
	}
	
	
	if isDebug{
		fmt.Println(wRusher, threadContralChan)
	}
}

func parseLine(line string, wRusher webRusher){
	// work finished, fill the main process chan and the thread contral chan
	fCloseMainChan := func(){mainProcessChan<-true}
	fCloseThreadContralChan := func(){<-threadContralChan}
	defer fCloseMainChan()
	if nil != threadContralChan{
		defer fCloseThreadContralChan()
	}
	
	isContain := strings.Contains(line, "MONITOR")
	if isContain{
		firstIdx := strings.Index(line, wRusher.startMark)
		if -1 == firstIdx{
			return
		}
		
		firstIdx += len(wRusher.startMark) + 1
		
		halfLine := line[firstIdx:]
		
		secondIdx := strings.Index(halfLine, wRusher.endMark)
		if -1 == secondIdx{
			return
		}
		
		targetLine := halfLine[0:secondIdx - 1]
		
		if isDebug{
			fmt.Println(targetLine)
		}
		
		var webParams map[string]interface{}
		err := json.Unmarshal([]byte(targetLine), &webParams)
		if nil != err{
			log.Println("failed to unmarshal json:", err, "targetLine:", targetLine)
			return
		}
		
		params := url.Values{}
		
		var strValue string
		for key, value := range webParams{
			switch targetValue := value.(type){
				case int:
				strValue = strconv.Itoa(targetValue)
				case string:
				strValue = targetValue
				default:
				continue
			}
			
			params.Add(key, strValue)
		}
		
		// 进行网络请求
		resp, err := http.PostForm(wRusher.targetUrl,params)
	
		if err != nil {
			log.Fatalln(err)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			return
		}

		if isDebug{
			fmt.Println(string(body))
		}
		
	}else{
		if isDebug{
			fmt.Println("not contain")
		}
	}
}

func (wRusher webRusher) DealWithLine(line string) {
	if nil != threadContralChan{
		threadContralChan<-true
	}
	
	if isDebug {
			fmt.Println("processing: ", workerNum)
		}
	
	workerNum++
	go parseLine(line, wRusher)
}

func (wRusher webRusher) dealWithReplayFile(){
	errno := fileReader.ReadLine(wRusher.replayFilePath, wRusher)
	if 0 != errno{
		log.Fatal("err occured while dealing with replay file")
	}
}

func Run(){	
	wRusher := new(webRusher)
	wRusher.getInitConf()		
	wRusher.dealWithReplayFile()
	
	wa := webAttacker.NewWebAttacker()
	fmt.Println(wa)
	
	// sleep some time for workers to start
	time.Sleep(time.Second * 1)
	
	var workerCount int64
	workerCount = 0
	for ;workerCount < workerNum; workerCount++{
		<-mainProcessChan
	}
	
	log.Println("all works done")
}