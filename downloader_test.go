package main

import "testing"

func TestAbbr(t *testing.T) {
	result := abbr(1024)

	if result != "1024" {
		t.Errorf("Result: %s", result)
	}
}
