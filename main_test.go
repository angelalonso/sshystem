package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	body := "ok"
	if body != "ok" {
		t.Errorf("Expected an 'ok' message. Got %s", body)
	}
}

//What should the script do?

// TEST: read config file
func TestReadConfig(t *testing.T) {
	testConns := readConfig("./machine_test.list")
	if len(testConns) != 2 {
		t.Errorf("Expected getting two objects. Got %d", len(testConns))
	}
	if testConns[0].User != "admin" {
		t.Errorf("Expected getting %s. Got %s", "admin", testConns[0].User)
	}
	if testConns[0].Host != "127.0.0.1" {
		t.Errorf("Expected getting %s. Got %s", "127.0.0.1", testConns[0].Host)
	}
	if testConns[0].Port != "22" {
		t.Errorf("Expected getting %s. Got %s", "22", testConns[0].Port)
	}

}

// TEST: ssh to machine
func TestSsh(t *testing.T) {
	expected_out := ""
	expected_err := "ssh: connect to host localhost port 22: Connection refused\r\n"
	out, err := sshCommand("localhost", "22", "pwd")
	if err != expected_err {
		t.Errorf("Expected getting error " + expected_err + " Got " + err)
	}
	if out != expected_out {
		t.Errorf("Expected getting " + expected_out + " Got " + out)
	}
}

// TODO: run command
// TODO: modify result
func TestFormatMem(t *testing.T) {
	out, err := mockCommand("/usr/bin/free", "")
	if err != "" {
		t.Errorf("Expected Not getting any error, but got " + err)
	}
	formatted := strings.Split(out, "\n")
	total, _ := strconv.Atoi(strings.Fields(formatted[1])[1])
	free, _ := strconv.Atoi(strings.Fields(formatted[1])[3])
	used := total - free
	percentage := float64(used) / float64(total) * 100
	fmt.Println(total)
	fmt.Println(free)
	fmt.Println(used)
	fmt.Printf("%0.2f %%", percentage)

}

func mockCommand(command string, params ...string) (string, string) {
	cmd := exec.Command(command, params...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	return outb.String(), errb.String()
}
