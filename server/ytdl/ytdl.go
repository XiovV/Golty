package ytdl

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

type Ytdl struct {
	baseCommand string
	resolutions map[string]string
	logger      *zap.Logger
}

func New(baseCommand string, logger *zap.Logger) *Ytdl {
	resolutions := map[string]string{
		"2160p": "2160",
		"1440p": "1440",
		"1080p": "1080",
		"720p":  "720",
		"480p":  "480",
		"360p":  "360",
		"240p":  "240",
		"144p":  "144",
	}

	return &Ytdl{baseCommand: baseCommand, resolutions: resolutions, logger: logger}
}

type VideoDownloadOptions struct {
	Video      bool
	Audio      bool
	Resolution string
	Format     string
	Output     string
}

func (y *Ytdl) exec(args ...string) ([]byte, error) {
	cmd := exec.Command(y.baseCommand, args...)

	ytdlOutput, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%s error: %s", y.baseCommand, exitError.Stderr)
		}

		return nil, fmt.Errorf("%s error: %w", y.baseCommand, err)
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
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s error: %s", "jq", exitError.Stderr)
		}

		return fmt.Errorf("%s error: %w", "jq", err)
	}

	return nil
}

func (y *Ytdl) downloadVideo(videoId string, options VideoDownloadOptions) error {
	if options.Video {
		resolutionFlag := fmt.Sprintf("res:%s", y.resolutions[options.Resolution])
		_, err := y.exec("-S", resolutionFlag, "-o", options.Output, videoId)
		return err
	}

	return nil
}
