package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	body := "ok"
	if body != "ok" {
		t.Errorf("Expected an 'ok' message. Got %s", body)
	}
}

//What should the script do?
// read config file
// ssh to machine
// run command
// get result
// modify result

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
