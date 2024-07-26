package plugins

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func FileRead(userFile string, passFile string) (userList []string, passList []string) {
	readFile := func(fileName string) ([]string, error) {
		var lines []string
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text()) // 去除换行符和左右空格
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return lines, nil
	}

	userList, err := readFile(userFile)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", userFile, err)
	}

	passList, err = readFile(passFile)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", passFile, err)
	}

	return userList, passList
}
