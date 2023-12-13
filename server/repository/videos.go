package repository

type Video struct {
	ID           int    `db:"id"`
	ChannelId    int    `db:"channelId"`
	VideoId      string `db:"videoId"`
	Title        string `db:"title"`
	ThumbnailUrl string `db:"thumbnailUrl"`
}

func (r *Repository) InsertVideo(video Video) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "INSERT INTO videos (channelId, videoId, title, thumbnailUrl) VALUES ($1, $2, $3, $4)", video.ChannelId, video.VideoId, video.Title, video.ThumbnailUrl)
	if err != nil {
		return err
	}

	return nil
}
