package config

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Hosts          []string
	User, Password string
	Port, Parallel int
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		profile_start, hosts_start bool
		hosts                      = make([]string, 0, 10)
		conf                       Config
	)

	r := bufio.NewReader(f)
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF && len(s) <= 0 {
			break
		}
		s = strings.Trim(s, "\n")

		switch {
		case strings.Contains(s, "profile"):
			profile_start = true
			continue

		case strings.Contains(s, "hosts"):
			profile_start, hosts_start = false, true
			continue

		case strings.Trim(s, " ") == "\n":
			continue
		}

		if profile_start {
			res := strings.Split(s, "=")
			switch strings.Trim(res[0], " ") {
			case "username":
				conf.User = strings.Trim(res[1], " ")

			case "password":
				conf.Password = strings.Trim(res[1], " ")

			case "port":
				conf.Port, _ = strconv.Atoi(strings.Trim(res[1], " "))

			case "parallel":
				conf.Parallel, _ = strconv.Atoi(strings.Trim(res[1], " "))
			}
		}

		if hosts_start {
			if strings.Contains(s, "-") {
				ips := strings.Split(s, "-")
				ip := strings.Split(ips[0], ".")

				start, err := strconv.Atoi(ip[3])
				if err != nil {
					log.Fatal("convert to integer fail: ", err)
				}

				stop, err := strconv.Atoi(ips[1])
				if err != nil {
					log.Fatal("convert to integer fail: ", err)
				}

				for i := start; i <= stop; i++ {
					ip[3] = strconv.Itoa(i)
					hosts = append(hosts, strings.Join(ip, "."))
				}
			} else {
				hosts = append(hosts, s)
			}
		}

	}
	conf.Hosts = hosts

	return &conf, nil
}
