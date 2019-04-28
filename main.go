package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Connection struct {
	User string
	Host string
	Port string
}

type Metric struct {
	Machine string
	Name    string
	Max     float32
	Current float32
}

func main() {
	//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	ex, _ := os.Executable()
	dir := filepath.Dir(ex)
	conns := readConfig(dir + "/machine.list")
	getResults(conns, dir)
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

func getResults(conns []Connection, maindir string) {
	// TODO: move the commands to a list file too
	metrics := []Metric{}
	for _, conn := range conns {
		outMem, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, "/usr/bin/free")
		metrics = append(metrics, getMetricMem(outMem, conn.Host))

		outTemp, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, "/opt/vc/bin/vcgencmd measure_temp")
		metrics = append(metrics, getMetricTemp(outTemp, conn.Host))

		outDisk, _ := sshCommand(conn.User+"@"+conn.Host, conn.Port, "df /")
		metrics = append(metrics, getMetricDisk(outDisk, conn.Host))
	}
	//fmt.Println(metrics)
	saveResults(metrics, maindir+"/metrics.csv")
}

func getMetricMem(memRaw string, machine string) Metric {
	formatted := strings.Split(memRaw, "\n")
	total, _ := strconv.Atoi(strings.Fields(formatted[1])[1])
	free, _ := strconv.Atoi(strings.Fields(formatted[1])[3])
	used := total - free
	metricMem := Metric{
		Machine: machine,
		Name:    "Mem",
		Max:     float32(total),
		Current: float32(used),
	}
	return metricMem
}

func getMetricTemp(tempRaw string, machine string) Metric {
	tempNDegrees := strings.Split(tempRaw, "=")[1]
	current64, _ := strconv.ParseFloat(tempNDegrees[:len(tempNDegrees)-3], 32)
	current := float32(current64)
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
		Max:     float32(total),
		Current: float32(used),
	}
	return metricDisk
}

func getPercentage(m Metric) float32 {
	percentage := float32(m.Current) / float32(m.Max) * 100
	return percentage
}

func saveResults(metrics []Metric, filename string) {
	t := time.Now()
	timestamp := fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, m := range metrics {
		entry := timestamp + ";" + m.Machine + ";" + m.Name + ";" + fmt.Sprintf("%f", m.Max) + ";" + fmt.Sprintf("%f", m.Current) + ";\n"
		if _, err = f.WriteString(entry); err != nil {
			panic(err)
		}
	}
}
