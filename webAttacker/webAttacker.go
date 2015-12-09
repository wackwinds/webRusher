// webAttacker project webAttacker.go
package webAttacker

import(
	"configReader"
	"strings"
	"net/http"
	"net/url"
	"log"
	"io/ioutil"
	"strconv"
	"encoding/json"
)

var ok bool

type WebAttacker struct {
	targetUrl string
	startMark string
	endMark string
}

const (
	INIT_CONF_NAME = "attackerConf.txt"
	IS_DEBUG = false
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
}

func (wa *WebAttacker) Attack(lineData interface{}){
	defer func(){
		if r := recover(); r != nil{
			log.Println("err in attack - reason:", r, " input:", lineData)
		}
	}()
	
	line := lineData.(string)
	isContain := strings.Contains(line, "MONITOR")
	if isContain{
		firstIdx := strings.Index(line, wa.startMark)
		if -1 == firstIdx{
			return
		}
		
		firstIdx += len(wa.startMark) + 1
		
		halfLine := line[firstIdx:]
		
		secondIdx := strings.Index(halfLine, wa.endMark)
		if -1 == secondIdx{
			return
		}
		
		targetLine := halfLine[0:secondIdx - 1]
		
		if IS_DEBUG{
			log.Println(targetLine)
		}
		
		var decordeResult interface{}
		err := json.Unmarshal([]byte(targetLine), &decordeResult)
		if nil != err{
			log.Println("failed to unmarshal json:", err, "targetLine:", targetLine)
			return
		}
		
		webParams := decordeResult.(map[string]interface{})
		
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
		resp, err := http.PostForm(wa.targetUrl,params)
	
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

		if IS_DEBUG{
			log.Println(string(body))
		}
		
	}else{
		if IS_DEBUG{
			log.Println("not contain")
		}
	}
}