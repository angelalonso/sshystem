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

func TestSsh(t *testing.T) {
	sshCommand("localhost", "22", "pwd")

}
