package utils

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"time"
)

func LocalCmd(cmd string) (string, error) {
	c := exec.Command("bash", "-c", cmd)
	output, err := c.CombinedOutput()
	return string(output), err
}

func LocalCmds(cmds []string) (string, error) {
	var outputs []string
	for _, cmd := range cmds {
		output, err := LocalCmd(cmd)
		if err != nil {
			return output, err
		}
		outputs = append(outputs, output)
	}
	return strings.Join(outputs, "\n"), nil
}

func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func sshClient(host, user, keyPath string) (*ssh.Client, error) {
	//key, err := ioutil.ReadFile("~/.ssh/id_rsa.pub")
	//if err != nil {
	//	return nil, err
	//}
	//signer, err := ssh.ParsePrivateKey(key)
	//if err != nil {
	//	return nil, err
	//}
	if keyPath == "" {
		keyPath = "/root/.ssh/id_rsa"
	}
	config := &ssh.ClientConfig{
		Timeout: 30 * time.Second,
		User:    user,
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: ssh.HostKeyCallback(func(string, net.Addr, ssh.PublicKey) error { return nil }),
		Auth:            []ssh.AuthMethod{publicKey(keyPath)},
	}
	c, err := ssh.Dial("tcp", host, config)
	return c, err
}

func SshCommand(host, user, keyPath, cmd string) (string, error) {
	c, err := sshClient(host, user, keyPath)
	if err != nil {
		return "", err
	}
	defer c.Close()

	session, err := c.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	return string(combo), nil
}

func SshCommands(host, user, keyPath string, cmds []string) (string, error) {
	c, err := sshClient(host, user, keyPath)
	if err != nil {
		return "", err
	}
	defer c.Close()

	session, err := c.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	combo, err := session.CombinedOutput(strings.Join(cmds, ";"))
	if err != nil {
		return "", err
	}
	return string(combo), nil
}
