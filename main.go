package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Connection struct {
	User string
	Host string
	Port string
}

type Metric struct {
	Name    string
	Machine string
	Max     float64
	Current float64
}

func main() {
	conns := readConfig("./machine.list")
	getResults(conns)
}

func readConfig(configfile string) []Connection {
	var connections []Connection
	file, err := os.Open(configfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) > 2 {
			newConnection := Connection{
				User: line[0],
				Host: line[1],
				Port: line[2],
			}
			connections = append(connections, newConnection)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return connections
}

func sshCommand(endpoint string, port string, command string) (string, string) {
	ssh_binary := "/usr/bin/ssh"
	cmd := exec.Command(ssh_binary, endpoint, "-p "+port, command)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	return outb.String(), errb.String()
}

// TODO: format this to receive result and work with it
func formatMem(memRaw string) {
	fmt.Println(memRaw)
}

func getResults(conns []Connection) {
	// TODO: move the commands to a list file too
	commands := []string{
		"hostname",
		"/opt/vc/bin/vcgencmd measure_temp",
		"df -h",
		"df -i"}
	for _, conn := range conns {
		for _, command := range commands {
			out, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, command)
			fmt.Println(out)
		}
	}
}
