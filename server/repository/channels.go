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
	ChannelId      int       `db:"channelId"`
	Quality        string    `db:"quality"`
	Format         string    `db:"format"`
	DownloadVideo  BoolAsInt `db:"downloadVideo"`
	DownloadAudio  BoolAsInt `db:"downloadAudio"`
	DownloadEntire BoolAsInt `db:"downloadEntire"`
	Sync           BoolAsInt `db:"sync"`
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

	_, err := r.db.ExecContext(ctx, "INSERT INTO channel_settings (channelId, quality, format, downloadVideo, downloadAudio, downloadEntire, sync) VALUES ($1, $2, $3, $4, $5, $6, $7)", settings.ChannelId, settings.Quality, settings.Format, settings.DownloadVideo, settings.DownloadAudio, settings.DownloadEntire, settings.Sync)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetChannels() ([]Channel, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channels := []Channel{}

	err := r.db.SelectContext(ctx, &channels, "SELECT id, channelName, channelHandle, channelUrl, avatarUrl ORDER BY id desc")
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (r *Repository) GetChannelsWithSize() ([]ChannelWithSize, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channels := []ChannelWithSize{}

	err := r.db.SelectContext(ctx, &channels, "SELECT channels.id, channels.channelName, channels.channelHandle, channels.channelUrl, channels.avatarUrl, COUNT(videos.videoId) as totalVideos, COALESCE(SUM(videos.size), 0) as totalSize FROM channels LEFT JOIN videos ON channels.id = videos.channelId GROUP BY channels.id ORDER BY channels.id desc")
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
	err := r.db.SelectContext(ctx, &channelDownloadSettings, "SELECT channelId, quality, format, downloadVideo, downloadAudio, downloadEntire, sync FROM channel_settings")
	if err != nil {
		return []ChannelDownloadSettings{}, err
	}

	return channelDownloadSettings, nil
}

func (r *Repository) GetChannelDownloadSettings(channelId int) (ChannelDownloadSettings, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channelDownloadSettings := ChannelDownloadSettings{}
	err := r.db.GetContext(ctx, &channelDownloadSettings, "SELECT quality, format, downloadVideo, downloadAudio, downloadEntire, sync FROM channel_settings WHERE channelId = $1", channelId)
	if err != nil {
		return ChannelDownloadSettings{}, err
	}

	return channelDownloadSettings, nil
}

func (r *Repository) DeleteChannel(channelId int) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM channels WHERE id = $1", channelId)
	if err != nil {
		return err
	}

	return err
}
