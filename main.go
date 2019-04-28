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
	Machine string
	Name    string
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
	metrics := []Metric{}
	for _, conn := range conns {
		outMem, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, "/usr/bin/free")
		//memMetric := getMetricMem(outMem, conn.Host)
		metrics = append(metrics, getMetricMem(outMem, conn.Host))

		outTemp, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, "/opt/vc/bin/vcgencmd measure_temp")
		//tempMetric := getMetricTemp(outTemp, conn.Host)
		metrics = append(metrics, getMetricTemp(outTemp, conn.Host))

		outDisk, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, "df /")
		//diskMetric := getMetricDisk(outDisk, conn.Host)
		metrics = append(metrics, getMetricDisk(outDisk, conn.Host))
	}
	fmt.Println(metrics)
}

func getMetricMem(memRaw string, machine string) Metric {
	formatted := strings.Split(memRaw, "\n")
	total, _ := strconv.Atoi(strings.Fields(formatted[1])[1])
	free, _ := strconv.Atoi(strings.Fields(formatted[1])[3])
	used := total - free
	metricMem := Metric{
		Machine: machine,
		Name:    "Mem",
		Max:     float64(total),
		Current: float64(used),
	}
	return metricMem
}

func getMetricTemp(tempRaw string, machine string) Metric {
	tempNDegrees := strings.Split(tempRaw, "=")[1]
	current, _ := strconv.ParseFloat(tempNDegrees[:len(tempNDegrees)-3], 64)
	metricTemp := Metric{
		Machine: machine,
		Name:    "Temp",
		Max:     0,
		Current: current,
	}
	return metricTemp
}

func getMetricDisk(diskRaw string, machine string) Metric {
	formatted := strings.Split(diskRaw, "\n")
	total, _ := strconv.Atoi(strings.Fields(formatted[1])[1])
	used, _ := strconv.Atoi(strings.Fields(formatted[1])[2])
	metricDisk := Metric{
		Machine: machine,
		Name:    "Disk",
		Max:     float64(total),
		Current: float64(used),
	}
	return metricDisk
}

func getPercentage(m Metric) float64 {
	percentage := float64(m.Current) / float64(m.Max) * 100
	return percentage
}
