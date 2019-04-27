package main

import (
	"bytes"
	"fmt"
	"os/exec"
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
	mem := formatMem(out)
	percentage := float64(mem.Current) / float64(mem.Max) * 100
	fmt.Printf("%0.2f %%\n", percentage)
}

func TestFormatTemp(t *testing.T) {
	out := "temp=52.1'C\n"
	temp := formatTemp(out)
	fmt.Printf("%0.2f °C\n", temp.Current)
}

func mockCommand(command string, params ...string) (string, string) {
	cmd := exec.Command(command, params...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	return outb.String(), errb.String()
}
