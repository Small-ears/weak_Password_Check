package plugins

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/malfunkt/iprange"
	"golang.com/golang.com/weak_Password_Check/vars"
)

// GET_IpList将u参数传入的ip地址范围进行解析并返回，得到一个ip，port结构体
func GET_IpList(ip string, port int) ([]vars.IpAddr, error) {
	add_List, err := iprange.ParseList(ip)
	if err != nil {
		return nil, err
	}

	var ipList []vars.IpAddr

	for _, addr := range add_List.Expand() {
		ipList = append(ipList, vars.IpAddr{
			Ip:   addr,
			Port: port,
		})
	}

	return ipList, nil
}

// GET_File_IpList读取传入的targetfile文件中的IP以及端口并进行分别存储
func GET_File_IpList(filePath string) ([]vars.IpAddr, error) {
	re := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	var ipList []vars.IpAddr
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") {
			// 将IP和端口存到结构体切片中
			parts := strings.Split(line, ":")
			ip := net.ParseIP(parts[0])
			port, err := strconv.Atoi(parts[1])
			if err != nil { //万一传入的不是数字而是字符
				log.Printf("Invalid port format: %s", parts[1])
				continue
			}
			if len(parts) == 2 {
				ipList = append(ipList, vars.IpAddr{Ip: ip, Port: port})
			}
		} else {
			if re.MatchString(line) {
				ipList = append(ipList, vars.IpAddr{Ip: net.ParseIP(line), Port: 3306})
			} else {
				log.Printf("Invalid IP format found: %s", line)
			}
		}
	}

	// 检查 bufio.Scanner 在扫描整个文件过程中的错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return ipList, err
}
