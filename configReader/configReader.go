// configReader project configReader.go
package configReader

import (
	"bufio"
	"fileReader"
	"fmt"
	"os"
	"strings"
)

var selfConfPath string

type ConfigReader struct {
	filePath  string
	ConfigMap map[string]string
}

func (cReader ConfigReader) UpdateConfFile() {
	file, err := os.Create(selfConfPath)
	if err != nil {
		fmt.Println("create file: " + selfConfPath + "failed")
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	for k, v := range cReader.ConfigMap {
		fmt.Fprintln(w, k+" : "+v)
	}
	w.Flush()
	return
}

func NewConfigReader(confPath string) *ConfigReader {
	selfConfPath = confPath
	configReader := new(ConfigReader)
	configReader.ConfigMap = make(map[string]string)

	fileReader.ReadLine(confPath, configReader)

	return configReader
}

func (cReader ConfigReader) GetConfig(confKey string) (string, bool) {
	val, ok := cReader.ConfigMap[confKey]
	return val, ok
}

func (cReader ConfigReader) ShowConfigMap() {
	fmt.Println(cReader.ConfigMap)
}

func (cReader ConfigReader) DealWithLine(line string) {
	trimLine := strings.TrimSpace(line)
	idx := strings.Index(trimLine, ":")
	if idx >= 0 {
		key := strings.TrimSpace(trimLine[0:idx])
		value := strings.TrimSpace(trimLine[idx+1:])
		cReader.ConfigMap[key] = value
	}
}
