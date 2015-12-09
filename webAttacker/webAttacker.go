// webAttacker project webAttacker.go
package webAttacker

import(
	"configReader"
	"log"
)

var ok bool

type WebAttacker struct {
	targetUrl string
	replayFilePath string
	startMark string
	endMark string
}

const (
	INIT_CONF_NAME = "attackerConf.txt"
)

func NewWebAttacker ()*WebAttacker{
	wa := new(WebAttacker)
	wa.init()
	return wa 
}

func (wa *WebAttacker) init(){
	conf := configReader.NewConfigReader(INIT_CONF_NAME)
	wa.targetUrl, ok = conf.GetConfig("targetUrl")
	if (!ok){
		log.Panic("failed to get target url")
	}
	
	wa.startMark, ok = conf.GetConfig("startMark")
	if (!ok){
		log.Panic("failed to get startMark")
	}
	
	wa.endMark, ok = conf.GetConfig("endMark")
	if (!ok){
		log.Panic("failed to get endMark")
	}
	
	wa.replayFilePath, ok = conf.GetConfig("replayFilePath")
	if (!ok){
		log.Panic("failed to get replay file path")
	}
}