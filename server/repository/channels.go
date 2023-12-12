package repository

type Channel struct {
	ID            int    `db:"id"`
	ChannelUrl    string `db:"channelUrl"`
	ChannelName   string `db:"channelName"`
	ChannelHandle string `db:"channelHandle"`
	AvatarUrl     string `db:"avatarUrl"`
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

func (r *Repository) GetChannels() ([]Channel, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channels := []Channel{}

	err := r.db.SelectContext(ctx, &channels, "SELECT id, channelUrl, channelName, channelHandle, avatarUrl FROM channels ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (r *Repository) GetChannelByHandle(channelHandle string) (Channel, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	channel := Channel{}

	err := r.db.GetContext(ctx, &channel, "SELECT id, channelUrl, channelName, channelHandle, avatarUrl FROM channels WHERE channelHandle = $1", channelHandle)
	if err != nil {
		return Channel{}, err
	}

	return channel, nil
}
