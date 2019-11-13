package main

import "testing"

func TestGetChannelName(t *testing.T) {
	result := GetChannelName("https://www.youtube.com/channel/UC4w1YQAJMWOz4qtxinq55LQ")

	if result != "UC4w1YQAJMWOz4qtxinq55LQ" {
		t.Errorf("Result Incorrect: %s", result)
	}
}
