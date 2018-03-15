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
	IdentityFile, IdentityPass string
	Port, Parallel int
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		profile, host 			   bool
		hosts                      = make([]string, 0, 10)
		conf                       Config
	)

	r := bufio.NewReader(f)
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF && len(s) <= 0 {
			break
		}
		s = strings.TrimSpace(s)

		switch {
		case strings.Contains(s, "[profile]"):
			profile = true
			continue

		case strings.Contains(s, "[hosts]"):
			profile, host = false, true
			continue

		case s == "" || strings.HasPrefix(s, "#"):
			continue
		}

		if profile {
			var value string
			res := strings.Split(s, "=")

			if len(res) == 2 {
				value = strings.TrimSpace(res[1])
			}

			switch strings.TrimSpace(res[0]) {
			case "username":
				conf.User = strings.Trim(value, `"`)

			case "password":
				conf.Password = strings.Trim(value,`"`)

			case "port":
				conf.Port, _ = strconv.Atoi(value)

			case "parallel":
				conf.Parallel, _ = strconv.Atoi(value)

			case "identityFile":
				conf.IdentityFile = strings.Trim(value, `"`)

			case "identityPass":
				conf.IdentityPass = strings.Trim(value, `"`)
			}
		}

		if host {
			if isDomainString(s) {
				hosts = append(hosts, s)
				continue
			}

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

func isDomainString(s string) bool {
	strList := strings.Split(s, ".")
	_, err := strconv.Atoi(strList[0])
	if err != nil {
		return true
	}
	return false
}