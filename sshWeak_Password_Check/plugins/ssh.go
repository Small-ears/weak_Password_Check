package plugins

import (
	"fmt"
	"net"
	"time"

	"golang.com/golang.com/sshWeak_Password_Check/vars"
	"golang.org/x/crypto/ssh"
)

func ScanSSH(s vars.Service) (result vars.ScanResult, err error) {
	result.Server = s

	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		Timeout: 5 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Target, s.Port), config)
	if err != nil {
		return result, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return result, err
	}

	err = session.Run("echo success") //一般情况下，如果err等于nil,则意味着 session.Run 成功执行了命令
	if err != nil {
		return result, err
	}

	result.Result = true
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
		if session != nil {
			_ = session.Close()
		}
	}()

	return result, err
}
