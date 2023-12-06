package ytdl

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type ChannelInfo struct {
	UploaderID string `json:"uploader_id"`
	Uploader   string `json:"uploader"`
	Avatar     struct {
		URL string `json:"url"`
	} `json:"avatar"`
}

func (y *Ytdl) GetChannelInfo(channelUrl string) (ChannelInfo, error) {
	ytdlArgs := []string{"--playlist-items", "0", "-J", channelUrl}

	cmd := exec.Command(y.BaseCommand, ytdlArgs...)

	ytdlOutput, err := cmd.Output()
	if err != nil {
		return ChannelInfo{}, err
	}

	jqArgs := []string{"{uploader_id, uploader, avatar: (.thumbnails[] | select(.id==\"avatar_uncropped\"))}"}

	jqCmd := exec.Command("jq", jqArgs...)

	jqCmd.Stdin = strings.NewReader(string(ytdlOutput))

	jqOutput, err := jqCmd.Output()
	if err != nil {
		return ChannelInfo{}, err
	}

	var channelInfo ChannelInfo
	err = json.Unmarshal(jqOutput, &channelInfo)
	if err != nil {
		return ChannelInfo{}, err
	}

	return channelInfo, err
}
