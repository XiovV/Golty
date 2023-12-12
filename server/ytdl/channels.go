package ytdl

import (
	"fmt"
)

type ChannelInfo struct {
	UploaderID string `json:"uploader_id"`
	Uploader   string `json:"uploader"`
	Avatar     struct {
		URL string `json:"url"`
	} `json:"avatar"`
}

type ChannelDownloadOptions struct {
	Video                    bool
	Audio                    bool
	Extension                string
	Quality                  string
	AutomaticallyDownloadNew bool
	DownloadEntire           bool
}

type ChannelVideos struct {
	Videos []string `json:"videos"`
}

func (y *Ytdl) GetChannelInfo(channelUrl string) (ChannelInfo, error) {
	ytdlOutput, err := y.exec("--playlist-items", "0", "-J", channelUrl)
	if err != nil {
		return ChannelInfo{}, err
	}

	var channelInfo ChannelInfo
	err = y.jq(string(ytdlOutput), &channelInfo, "{uploader_id, uploader, avatar: (.thumbnails[] | select(.id==\"avatar_uncropped\"))}")
	if err != nil {
		return ChannelInfo{}, err
	}

	return channelInfo, err
}

func (y *Ytdl) GetChannelVideos(channelUrl string) (ChannelVideos, error) {
	ytdlOutput, err := y.exec("--get-id", "--flat-playlist", channelUrl)
	if err != nil {
		return ChannelVideos{}, err
	}

	var channelVideos ChannelVideos
	err = y.jq(string(ytdlOutput), &channelVideos, "-nR", "{ \"videos\": [inputs] }")
	if err != nil {
		return ChannelVideos{}, err
	}

	return channelVideos, nil
}

func (y *Ytdl) DownloadChannel(channelUrl string, downloadOptions ChannelDownloadOptions) {
	channelVideos, err := y.GetChannelVideos(channelUrl)
	if err != nil {
		fmt.Println("could not get channel videos", err)
		return
	}

	fmt.Printf("%+v\n", channelVideos)
}
