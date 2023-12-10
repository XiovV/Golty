package repository

type Channel struct {
	ID            int
	ChannelUrl    string
	ChannelName   string
	ChannelHandle string
	AvatarUrl     string
}

func (r *Repository) InsertChannel(channel Channel) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "INSERT INTO channels (channelUrl, channelName, channelHandle, avatarUrl) VALUES ($1, $2, $3, $4)", channel.ChannelUrl, channel.ChannelName, channel.ChannelHandle, channel.AvatarUrl)
	if err != nil {
		return err
	}

	return nil
}
