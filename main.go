package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	//param_endpoint := "user@localhost"
	//param_port := "2"
	//param_command := `"whoami"`
	//sshCommand(param_endpoint, param_port, param_command)
	readConfig("./machine_test.list")
}

func readConfig(configfile string) {
	file, err := os.Open(configfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// TODO: create a struct, fill up an array of that struct and return it to be able to run ssh to each entry on that array
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func sshCommand(endpoint string, port string, command string) {
	ssh_binary := "/usr/bin/ssh"
	cmd := exec.Command(ssh_binary, endpoint, "-p "+port, command)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	fmt.Println("out:", outb.String(), "err:", errb.String(), err)
}
