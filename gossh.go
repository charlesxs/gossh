package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"gossh/config"
	"gossh/sshclient"
)

const example = `
Example <hosts.conf>:
[profile]
username = user
password = pass
identityFile = ~/.ssh/id_rsa
identityPass = "Null or private key's password"
port = 22
parallel = 10

[hosts]
192.168.1.3-20
192.168.1.100
192.168.1.253
`

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: ", os.Args[0], "command")
		return
	}

	var (
		wg    sync.WaitGroup
		count int
	)

	cmd := strings.Join(os.Args[1:], " ")
	conf, err := config.LoadConfig("hosts.conf")
	if err != nil {
		fmt.Println(err, example)
		return
	}

	for _, host := range conf.Hosts {
		count++
		wg.Add(1)
		go func(host string) {
			ssh := sshclient.NewSSH(
				host, conf.User,
				conf.Password,
				conf.IdentityFile,
				conf.IdentityPass,
				conf.Port)
			ssh.PrintRun(cmd)
			wg.Done()

		}(host)

		if count == conf.Parallel {
			wg.Wait()
			count = 0
		}
	}
	wg.Wait()
}
