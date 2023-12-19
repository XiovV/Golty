package service

import (
	"golty/repository"
	"golty/ytdl"
)

func (s *ChannelsService) channelOptionsToVideoOptions(channelOptions ChannelDownloadOptions, output string) ytdl.VideoDownloadOptions {
	return ytdl.VideoDownloadOptions{
		Video:      channelOptions.Video,
		Audio:      channelOptions.Audio,
		Resolution: channelOptions.Resolution,
		Output:     output,
	}
}

func (s *ChannelsService) channelOptionsToDBChannelOptions(channelId int, channelOptions ChannelDownloadOptions) repository.ChannelDownloadSettings {
	return repository.ChannelDownloadSettings{
		ChannelId:          channelId,
		Resolution:         channelOptions.Resolution,
		Format:             channelOptions.Format,
		DownloadVideo:      s.booleanToInteger(channelOptions.Video),
		DownloadAudio:      s.booleanToInteger(channelOptions.Audio),
		DownloadEntire:     s.booleanToInteger(channelOptions.DownloadEntire),
		DownloadNewUploads: s.booleanToInteger(channelOptions.DownloadNewUploads),
	}
}
