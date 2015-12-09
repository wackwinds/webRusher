// webRusher project webRusher.go
package webRusher

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"strconv"
	"configReader"
	"fileReader"
	"time"
	"webAttacker"
)

const (
	INIT_CONF_NAME = "rusherConf.txt"
	DEFAULTE_CPU_NUM = 0
	DEFAULTE_THREAD_NUM = 0
)

var isDebug bool
var ok bool
var workerMutex sync.Mutex
var workerNum int64
var mainProcessChan chan bool
var threadContralChan chan bool

type webRusher struct {
	attacker webAttacker.WebAttacker
	replayFilePath string
}

func init(){
	isDebug = false
	workerNum = 0;
	mainProcessChan = make(chan bool)
}

func (wr *webRusher) getInitConf(){
	conf := configReader.NewConfigReader(INIT_CONF_NAME)
	
	wr.replayFilePath, ok = conf.GetConfig("replayFilePath")
	if (!ok){
		log.Panic("failed to get replay file path")
	}
	
	strCpuNumForThread, ok := conf.GetConfig("cpuNumForThread")
	
	cpuNumForThread, err := strconv.Atoi(strCpuNumForThread)
	if nil != err{
		log.Panicln(err)
		cpuNumForThread = DEFAULTE_CPU_NUM
	}
	
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
	
	numOfThread, err := strconv.Atoi(strNumOfThread)
	if nil != err{
		log.Println(err)
		numOfThread = DEFAULTE_THREAD_NUM
	}
	
	if (!ok || (numOfThread <= 0)){
		log.Println("failed to get threadNum")
		threadContralChan = nil
	}else{
		threadContralChan = make(chan bool, numOfThread)
	}
}

func parseLine(line string, wr webRusher){
	// work finished, fill the main process chan and the thread contral chan
	defer func(){mainProcessChan<-true}()
	if nil != threadContralChan{
		defer func(){<-threadContralChan}()
	}
	
	wr.attacker.Attack(line)
}

func (wr webRusher) DealWithLine(line string) {
	if nil != threadContralChan{
		threadContralChan<-true
	}
	
	if isDebug {
			fmt.Println("processing: ", workerNum)
		}
	
	workerNum++
	go parseLine(line, wr)
}

func (wr webRusher) dealWithReplayFile(){
	errno := fileReader.ReadLine(wr.replayFilePath, wr)
	if 0 != errno{
		log.Fatal("err occured while dealing with replay file")
	}
}

func Run(){	
	wr := new(webRusher)
	wr.getInitConf()		
	wr.dealWithReplayFile()
	
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