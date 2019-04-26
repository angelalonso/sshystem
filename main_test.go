package main

import (
  "testing"
)

func TestMain(t *testing.T) {
  body := returnResult()
  if body != "ok" {
    t.Errorf("Expected an 'ok' message. Got %s", body)
  }
}
