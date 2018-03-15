package sshclient

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"time"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
)

type SSH struct {
	Host, User, Password string
	IdentityFile, IdentityFilePassword string
	Port                 int
}

func NewSSH(host, user, password, identityFile, idfilepass string, port int) *SSH {
	return &SSH{
		Host:     host,
		User:     user,
		Password: password,
		Port:     port,
		IdentityFile: identityFile,
		IdentityFilePassword: idfilepass,
	}
}

func (s *SSH) Connect() (*ssh.Session, error) {
	auths := []ssh.AuthMethod{
		ssh.Password(s.Password),
	}
	key, err := ioutil.ReadFile(s.IdentityFile)
	if err == nil {
		key, err := DecryptKey(key, []byte(s.IdentityFilePassword))
		if err != nil {
			goto M
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err == nil {
			auths = append(auths, ssh.PublicKeys(signer))
		}
	}
	M:
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: auths,
		Timeout: time.Second * 10,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return client.NewSession()
}

func (s *SSH) Run(cmd string) (stdout, stderr io.Reader, err error) {
	session, err := s.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	defer session.Close()

	stdout, _ = session.StdoutPipe()
	stderr, _ = session.StderrPipe()

	err = session.Run(cmd)
	return stdout, stderr, err
}

func (s *SSH) PrintRun(cmd string) int {
	var res []byte

	// initial session
	session, err := s.Connect()
	if err != nil {
		log.Printf("%s execute failed ~> \n\033[31m%s\033[0m\n\n", s.Host, err)
		return 1
	}
	defer session.Close()

	// set stdout, stderr
	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()

	// run command
	err = session.Run(cmd)
	if err != nil {
		res, _ = ioutil.ReadAll(stderr)
		log.Printf("%s execute failed ~> \n\033[31m%s\033[0m", s.Host, string(res))
		return 1
	}
	res, _ = ioutil.ReadAll(stdout)
	log.Printf("%s execute ok ~> \n\033[32m%s\033[0m", s.Host, string(res))
	return 0
}


func DecryptKey(key, password []byte) ([]byte, error) {
	block, rest := pem.Decode(key)
	if len(rest) > 0 {
		return nil, errors.New(fmt.Sprintf("Decrypt key error: %s", string(rest)))
	}

	if x509.IsEncryptedPEMBlock(block) {
		der, err := x509.DecryptPEMBlock(block, password)
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(&pem.Block{Type: block.Type, Bytes: der}), nil
	}
	return key, nil
}