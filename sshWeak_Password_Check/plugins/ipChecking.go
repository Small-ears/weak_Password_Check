package plugins

import (
	"fmt"
	"net"
	"sync"
	"time"

	"golang.com/golang.com/sshWeak_Password_Check/vars"
)

// Active_Checking最终获取的IP:port进行端口存活检测
func Active_Checking(addr []vars.IpAddr, maxConcurrency int) []vars.IpAddr {
	var wg sync.WaitGroup //确保所有 goroutines 完成后主程序继续执行
	var mu sync.Mutex
	semaphore := make(chan struct{}, maxConcurrency) //创建一个channel，channel的大小就是最大并发数量。
	var activeAddrs []vars.IpAddr

	for _, ipInfo := range addr {
		wg.Add(1)
		go func(ip net.IP, port int) {
			defer wg.Done() //计数器

			semaphore <- struct{}{}        //发送一个空结构体,占位符，如果信号量通道已满（例如，已有 10 个 goroutines 在运行），则这个操作会阻塞，直到有一个槽位被释放。
			defer func() { <-semaphore }() //匿名函数，确保在 goroutine 完成时释放一个信号量槽位

			address := fmt.Sprintf("%v:%v", ip, port)
			conn, err := net.DialTimeout("tcp", address, 3*time.Second)
			if err != nil {
				//fmt.Printf("Failed to connect to %v: %v\n", address, err)
				return
			}
			conn.Close()

			mu.Lock() //在访问共享资源newIpList时使用互斥锁，以确保并发安全。
			activeAddrs = append(activeAddrs, vars.IpAddr{Ip: ip, Port: port})
			mu.Unlock()
		}(ipInfo.Ip, ipInfo.Port)
	}

	wg.Wait()
	return activeAddrs
}
