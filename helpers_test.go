package main

import "testing"

func TestGetChannelName(t *testing.T) {
	result, err := GetChannelName("https://www.youtube.com/channel/UC4w1YQAJMWOz4qtxinq55LQ")

	if result != "UC4w1YQAJMWOz4qtxinq55LQ" {
		t.Errorf("Result Incorrect: %s", result)
	}
	if err != nil {
		t.Errorf("Error Set: %s", err)
	}
}
