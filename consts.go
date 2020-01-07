package main

const (
	CONFIG_ROOT = "./config/"
	VERSION     = "v0.5.11"
)

var (
	errRes Response
	okRes  Response
)

var checkingIntervalMultipliers = map[string]int{
	"minutes": 1,
	"hours":   60,
	"days":    1440,
}
