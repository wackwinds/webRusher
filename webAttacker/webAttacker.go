// webAttacker project webAttacker.go
package webAttacker

import(
	"configReader"
	"strings"
	"net/http"
	"net/url"
	"log"
	"runtime/debug"
	"io/ioutil"
	"strconv"
	"encoding/json"
)

const (
	INIT_CONF_NAME = "attackerConf.txt"
	IS_DEBUG = false
)

var ok bool

type WebAttacker struct {
	startMark string
	endMark string
	mode string
	host string
	port string
}

func NewWebAttacker ()*WebAttacker{
	wa := new(WebAttacker)
	wa.init()
	return wa 
}

func (wa *WebAttacker) init(){
	conf := configReader.NewConfigReader(INIT_CONF_NAME)
	
	wa.mode, ok = conf.GetConfig("mode")
	if (!ok){
		log.Println("failed to get attack mode")
	}
	
	wa.host, ok = conf.GetConfig("host")
	if (!ok){
		log.Panic("failed to get target host")
	}
	
	wa.port, ok = conf.GetConfig("port")
	if (!ok){
		log.Panic("failed to get target port")
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

func (wa *WebAttacker) BackDeal(respData interface{})(err error){
	if nil == respData{
		return nil
	}
	
	strData := string(respData.([]uint8))
	log.Println(strData)
	return nil
}

func (wa *WebAttacker) Attack(lineData interface{})(responsData interface{}, err error){
	defer func(){
		if r := recover(); r != nil{
			log.Println("err in attack - reason:", r, " input:", lineData)
			if IS_DEBUG{
				debug.PrintStack()
				log.Println("self:", wa)
			}
		}
	}()
	
	line := lineData.(string)
	isContain := strings.Contains(line, "MONITOR")
	
	if isContain{
		firstIdx := strings.Index(line, wa.startMark)
		if -1 == firstIdx{
			return nil, nil
		}
		
		uri := getUriInLine(line)
		
		firstIdx += len(wa.startMark) + 1
		
		halfLine := line[firstIdx:]
		
		secondIdx := strings.Index(halfLine, wa.endMark)
		if -1 == secondIdx{
			return nil, nil
		}
		
		if IS_DEBUG{
			log.Println("half line:", halfLine)
		}
		
		targetLine := halfLine[0:secondIdx - 1]
		
		if IS_DEBUG{
			log.Println("target line:", targetLine)
		}
		
		var decordeResult interface{}
		err := json.Unmarshal([]byte(targetLine), &decordeResult)
		if nil != err{
			log.Println("failed to unmarshal json:", err, "targetLine:", targetLine)
			return nil, err
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
		
		targetUrl := wa.host + ":" + wa.port + "/" + uri
		
		// 进行网络请求
		resp, err := http.PostForm(targetUrl,params)
	
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			return nil, err
		}

		if IS_DEBUG{
			log.Println(string(body))
		}
		
		return body, nil
		
	}else{
		if IS_DEBUG{
			log.Println("not contain")
		}
		
		return nil, nil
	}
}

/*
	在百度ODP日志里取出uri
*/
func getUriInLine(line string)(uri string){
	startMark := "uri["
	endMark := "]"
	firstIdx := strings.Index(line, startMark)
	
	firstIdx += len(startMark) + 1
	halfLine := line[firstIdx:]
	secondIdx := strings.Index(halfLine, endMark)
	
	uri = halfLine[:secondIdx]
	return
}