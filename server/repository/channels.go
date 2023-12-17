package repository

type Channel struct {
	ID            int    `db:"id"`
	ChannelUrl    string `db:"channelUrl"`
	ChannelName   string `db:"channelName"`
	ChannelHandle string `db:"channelHandle"`
	AvatarUrl     string `db:"avatarUrl"`
}

type ChannelWithSize struct {
	Channel
	TotalVideos int `db:"totalVideos"`
	TotalSize   int `db:"totalSize"`
}

type ChannelDownloadSettings struct {
	ChannelId          int    `db:"channelId"`
	Resolution         string `db:"resolution"`
	Format             string `db:"format"`
	DownloadVideo      int    `db:"downloadVideo"`
	DownloadAudio      int    `db:"downloadAudio"`
	DownloadEntire     int    `db:"downloadEntire"`
	DownloadNewUploads int    `db:"downloadNewUploads"`
}

func (r *Repository) InsertChannel(channel Channel) (Channel, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	var newChannel Channel
	err := r.db.GetContext(ctx, &newChannel, "INSERT INTO channels (channelUrl, channelName, channelHandle, avatarUrl) VALUES ($1, $2, $3, $4) RETURNING *;", channel.ChannelUrl, channel.ChannelName, channel.ChannelHandle, channel.AvatarUrl)
	if err != nil {
		return Channel{}, err
	}

	return newChannel, err
}

func (r *Repository) InsertChannelDownloadSettings(settings ChannelDownloadSettings) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "INSERT INTO channel_settings (channelId, resolution, format, downloadVideo, downloadAudio, downloadEntire, downloadNewUploads) VALUES ($1, $2, $3, $4, $5, $6, $7)", settings.ChannelId, settings.Resolution, settings.Format, settings.DownloadVideo, settings.DownloadAudio, settings.DownloadEntire, settings.DownloadNewUploads)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetChannels() ([]Channel, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channels := []Channel{}

	err := r.db.SelectContext(ctx, &channels, "SELECT id, channelName, channelHandle, channelUrl, avatarUrl")
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (r *Repository) GetChannelsWithSize() ([]ChannelWithSize, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channels := []ChannelWithSize{}

	err := r.db.SelectContext(ctx, &channels, "SELECT channels.id, channels.channelName, channels.channelHandle, channels.channelUrl, channels.avatarUrl, COUNT(videos.videoId) as totalVideos, COALESCE(SUM(videos.size), 0) as totalSize FROM channels LEFT JOIN videos ON channels.id = videos.channelId GROUP BY channels.id")
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (r *Repository) FindChannelByHandle(channelHandle string) (Channel, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channel := Channel{}

	err := r.db.GetContext(ctx, &channel, "SELECT id, channelName, channelHandle, channelUrl, avatarUrl FROM channels WHERE channelHandle = $1", channelHandle)
	if err != nil {
		return Channel{}, err
	}

	return channel, nil
}

func (r *Repository) FindChannelByID(channelId int) (Channel, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channel := Channel{}

	err := r.db.GetContext(ctx, &channel, "SELECT id, channelName, channelHandle, channelUrl, avatarUrl FROM channels WHERE id = $1", channelId)
	if err != nil {
		return Channel{}, err
	}

	return channel, nil
}

func (r *Repository) FindChannelByHandleWithSize(channelHandle string) (ChannelWithSize, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channel := ChannelWithSize{}

	err := r.db.GetContext(ctx, &channel, "SELECT channels.id, channels.channelName, channels.channelHandle, channels.channelUrl, channels.avatarUrl, COUNT(videos.videoId) as totalVideos, COALESCE(SUM(videos.size), 0) as totalSize FROM channels LEFT JOIN videos ON channels.id = videos.channelId WHERE channels.channelHandle = $1", channelHandle)
	if err != nil {
		return ChannelWithSize{}, err
	}

	return channel, nil
}

func (r *Repository) GetAllDownloadSettings() ([]ChannelDownloadSettings, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channelDownloadSettings := []ChannelDownloadSettings{}
	err := r.db.SelectContext(ctx, &channelDownloadSettings, "SELECT channelId, resolution, format, downloadVideo, downloadAudio, downloadEntire, downloadNewUploads FROM channel_settings")
	if err != nil {
		return []ChannelDownloadSettings{}, err
	}

	return channelDownloadSettings, nil
}
