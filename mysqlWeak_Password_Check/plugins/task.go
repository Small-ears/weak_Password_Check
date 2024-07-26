package plugins

import "golang.com/golang.com/weak_Password_Check/vars"

func GenerateTask(ipList []vars.IpAddr, users []string, passwords []string) (tasks []vars.Service, taskNum int) {
	tasks = make([]vars.Service, 0)

	for _, user := range users {
		for _, password := range passwords {
			for _, addr := range ipList {
				service := vars.Service{Target: addr.Ip, Port: addr.Port, Username: user, Password: password}
				tasks = append(tasks, service)
			}
		}
	}

	return tasks, len(tasks)
}
