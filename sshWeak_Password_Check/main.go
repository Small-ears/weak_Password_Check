package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"golang.com/golang.com/sshWeak_Password_Check/plugins"
	"golang.com/golang.com/sshWeak_Password_Check/vars"
)

func main() {
	var target, targetFile, userFile, passFile string
	var port, goroutineNum int
	flag.StringVar(&target, "u", "", "Enter target address, Example: 10.0.0.1,10.0.0.0/24,10.0.0.*,10.0.0.1-10")
	flag.IntVar(&port, "p", 22, "Enter the port number, default is 22")
	flag.StringVar(&targetFile, "r", "", "Enter target address list file, Example: path\\ips.txt,Supported formats:10.0.0.1 or 10.0.0.1:2222")
	flag.StringVar(&userFile, "U", "", "Enter username list file, Example: path\\userFile.txt")
	flag.StringVar(&passFile, "P", "", "Enter password list file, Example: path\\passFile.txt")
	flag.IntVar(&goroutineNum, "g", 10, "Enter the number of goroutines, default is 10")

	// 自定义Usage函数
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// 检查是否有参数传入
	if len(os.Args) < 7 {
		flag.Usage()
		os.Exit(1)
	}

	if target != "" && targetFile != "" {
		log.Fatal("-u和-r不能一起使用")
		flag.Usage()
		os.Exit(1)
	}

	userList, passList := plugins.FileRead(userFile, passFile)

	var activeAddrs []vars.IpAddr //在 if 块外没有被使用，这在 Go 语言中会被视为变量未使用而报错。

	if target != "" {
		uIpList, err := plugins.GET_IpList(target, port)
		if err != nil {
			log.Fatal("Error parsing the provided IP address.:", err)
		}
		activeAddrs = plugins.Active_Checking(uIpList, goroutineNum)
		fmt.Println("Completed target port scan,Starting password brute force attack.")
	}

	if targetFile != "" {
		fIpList, err := plugins.GET_File_IpList(targetFile)
		if err != nil {
			log.Fatalf("File read error")
		}
		activeAddrs = plugins.Active_Checking(fIpList, goroutineNum)
		fmt.Println("Completed target port scan,Starting password brute force attack.")
	}

	ipList, num := plugins.GenerateTask(activeAddrs, userList, passList)

	var wg sync.WaitGroup
	var m sync.Mutex
	wg.Add(num)
	workerChan := make(chan struct{}, goroutineNum) //控制并发大小
	for _, taskList := range ipList {
		workerChan <- struct{}{} //占用一个并发信号量
		go func(s vars.Service) {
			defer func() {
				<-workerChan // 释放一个并发信号量
				wg.Done()
			}()
			result, err := plugins.ScanSSH(s)
			if err != nil {
				return
			}

			fmt.Println(result.Server.Target, result.Server.Username, result, result)

			m.Lock() //当在并发中写动作时，需要使用互斥锁
			plugins.WriteScanResultToTXT(result)
			m.Unlock()
		}(taskList)
	}

	wg.Wait()
	fmt.Println("All scans completed.")
}
