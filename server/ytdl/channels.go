package ytdl

const CHANNELS_DEFAULT_OUTPUT = "videos/channels/%(uploader)s/%(id)s.%(ext)s"

type ChannelInfo struct {
	UploaderID string `json:"uploader_id"`
	Uploader   string `json:"uploader"`
	Avatar     struct {
		URL string `json:"url"`
	} `json:"avatar"`
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

func (y *Ytdl) GetChannelVideos(channelUrl string) ([]string, error) {
	ytdlOutput, err := y.exec("--get-id", "--flat-playlist", channelUrl)
	if err != nil {
		return nil, err
	}

	var channelVideos ChannelVideos
	err = y.jq(string(ytdlOutput), &channelVideos, "-nR", "{ \"videos\": [inputs] }")
	if err != nil {
		return nil, err
	}

	return channelVideos.Videos, nil
}
