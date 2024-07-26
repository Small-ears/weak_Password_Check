package plugins

import (
	"fmt"
	"os"

	"golang.com/golang.com/weak_Password_Check/vars"
)

func WriteScanResultToTXT(result vars.ScanResult) {
	file, err := os.OpenFile("Result.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Write error, result is:", result.Server.Target, result.Server.Username, result.Server.Password, result.Result)
	}
	defer file.Close()

	line := fmt.Sprintf("Target: %v, Username: %v, Password: %v, Result: %v\n",
		result.Server.Target, result.Server.Username, result.Server.Password, result.Result)

	_, err = file.WriteString(line)

	if err != nil {
		fmt.Println("Write error, result is:", result.Server.Target, result.Server.Username, result.Server.Password, result.Result)
	}
}
