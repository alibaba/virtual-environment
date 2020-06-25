package parser

import (
	"testing"
)

func TestGetKeys(t *testing.T) {
	keys := GetKeys(map[string]bool{"a": true, "c": false, "e": true})
	if len(keys) != 3 {
		t.Fail()
	}
	if keys[0] != "a" || keys[1] != "c" || keys[2] != "e" {
		t.Fail()
	}
}
