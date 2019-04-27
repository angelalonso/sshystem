package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
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

func getResults(conns []Connection) {
	// TODO: move the commands to a list file too
	commands := []string{
		"df -h",
		"df -i"}
	for _, conn := range conns {
		fmt.Println(conn.Host)

		outMem, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, "/usr/bin/free")
		memMetric := formatMem(outMem)
		percentage := float64(memMetric.Current) / float64(memMetric.Max) * 100
		fmt.Printf("Mem: %0.2f%%\n", percentage)

		outTemp, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, "/opt/vc/bin/vcgencmd measure_temp")
		tempMetric := formatTemp(outTemp)
		fmt.Printf("Temp: %0.2f °C\n", tempMetric.Current)

		for _, command := range commands {
			out, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, command)
			fmt.Println(out)
		}
	}
}

func formatMem(memRaw string) Metric {
	formatted := strings.Split(memRaw, "\n")
	total, _ := strconv.Atoi(strings.Fields(formatted[1])[1])
	free, _ := strconv.Atoi(strings.Fields(formatted[1])[3])
	used := total - free
	memMetric := Metric{
		Name:    "Mem",
		Machine: "",
		Max:     float64(total),
		Current: float64(used),
	}
	return memMetric
}

func formatTemp(tempRaw string) Metric {
	tempNDegrees := strings.Split(tempRaw, "=")[1]
	current, _ := strconv.ParseFloat(tempNDegrees[:len(tempNDegrees)-3], 64)
	memMetric := Metric{
		Name:    "Temp",
		Machine: "",
		Max:     0,
		Current: current,
	}
	return memMetric
}
