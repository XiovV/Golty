package repository

type Video struct {
	ID           int    `db:"id"`
	ChannelId    int    `db:"channelId"`
	VideoId      string `db:"videoId"`
	Title        string `db:"title"`
	ThumbnailUrl string `db:"thumbnailUrl"`
	Size         int64  `db:"size"`
	DownloadDate int64  `db:"downloadDate"`
	UploadDate   int64  `db:"uploadDate"`
}

func (r *Repository) InsertVideo(video Video) error {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "INSERT INTO videos (channelId, videoId, title, thumbnailUrl, size, downloadDate, uploadDate) VALUES ($1, $2, $3, $4, $5, $6, $7)", video.ChannelId, video.VideoId, video.Title, video.ThumbnailUrl, video.Size, video.DownloadDate, video.UploadDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetChannelVideos(channelId int) ([]Video, error) {
	ctx, cancel := newBackgroundContext(DefaultQueryTimeout)
	defer cancel()

	videos := []Video{}
	err := r.db.SelectContext(ctx, &videos, "SELECT id, channelId, videoId, title, thumbnailUrl, size, downloadDate, uploadDate FROM videos WHERE channelId = $1 ORDER BY uploadDate DESC", channelId)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
