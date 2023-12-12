package ytdl

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Ytdl struct {
	BaseCommand string
}

func New(baseCommand string) *Ytdl {
	return &Ytdl{BaseCommand: baseCommand}
}

func (y *Ytdl) exec(args ...string) ([]byte, error) {
	cmd := exec.Command(y.BaseCommand, args...)

	ytdlOutput, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%s error: %s", y.BaseCommand, exitError.Stderr)
		}

		return nil, fmt.Errorf("%s error: %w", y.BaseCommand, err)
	}

	return ytdlOutput, nil
}

func (y *Ytdl) jq(input string, output any, args ...string) error {
	jqCmd := exec.Command("jq", args...)

	jqCmd.Stdin = strings.NewReader(input)

	jqOutput, err := jqCmd.Output()
	if err != nil {
		return err
	}
	err = json.Unmarshal(jqOutput, output)
	if err != nil {
		return err
	}

	return nil
}
