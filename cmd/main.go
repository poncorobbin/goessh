package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var stdout, stderr bytes.Buffer

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <user> <host:port> <command>", os.Args[0])
	}
	// create ssh connection
	client, session := connect(os.Args[2], os.Args[1])
	if client != nil {
		defer client.Close()
		defer session.Close()
	}

	// get commands from args
	argsRaw := os.Args
	commands := argsRaw[3]

	out, err := session.CombinedOutput(commands)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n============== OUTPUT")
	fmt.Println(string(out))
}

func connect(sshAddress string, sshUser string) (*ssh.Client, *ssh.Session) {
	fmt.Print("Password: ")
	// fmt.Scanf("%s\n", &sshPwd)
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		os.Exit(1)
	}
	sshPwd := string(bytepw)

	sshConfig := &ssh.ClientConfig{
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPwd),
		},
	}

	client, err := ssh.Dial("tcp", sshAddress, sshConfig)
	if err != nil {
		log.Fatal("Failed to dial ", err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session ", err.Error())
	}

	return client, session
}
