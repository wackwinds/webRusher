// fileReader project fileReader.go
package fileReader

import (
	"bufio"
	"os"
)

type lineHandler interface {
	DealWithLine(string)
}

func ReadLine(filePath string, handler lineHandler) (errno int) {
	file, err := os.Open(filePath)
	if err != nil {
		return 1
	}
	defer file.Close()

	// var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		handler.DealWithLine(line)
	}
	return 0
}
