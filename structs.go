package main

type TargetMetadata struct {
	PlaylistUploader   string      `json:"playlist_uploader"`
	UploadDate         string      `json:"upload_date"`
	Extractor          string      `json:"extractor"`
	Series             interface{} `json:"series"`
	Format             string      `json:"format"`
	Vbr                interface{} `json:"vbr"`
	Chapters           interface{} `json:"chapters"`
	Height             int         `json:"height"`
	LikeCount          int         `json:"like_count"`
	Duration           int         `json:"duration"`
	Fulltitle          string      `json:"fulltitle"`
	PlaylistIndex      int         `json:"playlist_index"`
	Album              string      `json:"album"`
	ViewCount          int         `json:"view_count"`
	Playlist           string      `json:"playlist"`
	Title              string      `json:"title"`
	Filename           string      `json:"_filename"`
	Creator            string      `json:"creator"`
	Ext                string      `json:"ext"`
	ID                 string      `json:"id"`
	DislikeCount       int         `json:"dislike_count"`
	PlaylistID         string      `json:"playlist_id"`
	Abr                int         `json:"abr"`
	UploaderURL        string      `json:"uploader_url"`
	Categories         []string    `json:"categories"`
	Fps                int         `json:"fps"`
	StretchedRatio     interface{} `json:"stretched_ratio"`
	SeasonNumber       interface{} `json:"season_number"`
	Annotations        interface{} `json:"annotations"`
	WebpageURLBasename string      `json:"webpage_url_basename"`
	Acodec             string      `json:"acodec"`
	DisplayID          string      `json:"display_id"`
	RequestedFormats   []struct {
		Asr               interface{} `json:"asr"`
		Tbr               float64     `json:"tbr"`
		Protocol          string      `json:"protocol"`
		Format            string      `json:"format"`
		URL               string      `json:"url"`
		Vcodec            string      `json:"vcodec"`
		FormatNote        string      `json:"format_note"`
		Height            int         `json:"height"`
		DownloaderOptions struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options"`
		Width       int    `json:"width"`
		Ext         string `json:"ext"`
		Filesize    int    `json:"filesize"`
		Fps         int    `json:"fps"`
		FormatID    string `json:"format_id"`
		PlayerURL   string `json:"player_url"`
		HTTPHeaders struct {
			AcceptCharset  string `json:"Accept-Charset"`
			AcceptLanguage string `json:"Accept-Language"`
			AcceptEncoding string `json:"Accept-Encoding"`
			Accept         string `json:"Accept"`
			UserAgent      string `json:"User-Agent"`
		} `json:"http_headers"`
		Acodec string `json:"acodec"`
		Abr    int    `json:"abr,omitempty"`
	} `json:"requested_formats"`
	AutomaticCaptions struct {
	} `json:"automatic_captions"`
	Description        string      `json:"description"`
	Tags               []string    `json:"tags"`
	Track              string      `json:"track"`
	RequestedSubtitles interface{} `json:"requested_subtitles"`
	StartTime          interface{} `json:"start_time"`
	AverageRating      float64     `json:"average_rating"`
	Uploader           string      `json:"uploader"`
	FormatID           string      `json:"format_id"`
	EpisodeNumber      interface{} `json:"episode_number"`
	UploaderID         string      `json:"uploader_id"`
	Subtitles          struct {
	} `json:"subtitles"`
	PlaylistTitle string `json:"playlist_title"`
	Width         int    `json:"width"`
	Thumbnails    []struct {
		URL string `json:"url"`
		ID  string `json:"id"`
	} `json:"thumbnails"`
	License      interface{} `json:"license"`
	Artist       string      `json:"artist"`
	ExtractorKey string      `json:"extractor_key"`
	ReleaseDate  interface{} `json:"release_date"`
	AltTitle     string      `json:"alt_title"`
	Thumbnail    string      `json:"thumbnail"`
	ChannelID    string      `json:"channel_id"`
	IsLive       interface{} `json:"is_live"`
	ReleaseYear  interface{} `json:"release_year"`
	EndTime      interface{} `json:"end_time"`
	WebpageURL   string      `json:"webpage_url"`
	Formats      []struct {
		Asr               int     `json:"asr"`
		Tbr               float64 `json:"tbr"`
		Protocol          string  `json:"protocol"`
		Format            string  `json:"format"`
		URL               string  `json:"url"`
		Vcodec            string  `json:"vcodec"`
		FormatNote        string  `json:"format_note"`
		Abr               int     `json:"abr,omitempty"`
		PlayerURL         string  `json:"player_url"`
		DownloaderOptions struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options,omitempty"`
		Width       interface{} `json:"width"`
		Ext         string      `json:"ext"`
		Filesize    int         `json:"filesize"`
		Fps         interface{} `json:"fps"`
		FormatID    string      `json:"format_id"`
		Height      interface{} `json:"height"`
		HTTPHeaders struct {
			AcceptCharset  string `json:"Accept-Charset"`
			AcceptLanguage string `json:"Accept-Language"`
			AcceptEncoding string `json:"Accept-Encoding"`
			Accept         string `json:"Accept"`
			UserAgent      string `json:"User-Agent"`
		} `json:"http_headers"`
		Acodec    string `json:"acodec"`
		Container string `json:"container,omitempty"`
	} `json:"formats"`
	PlaylistUploaderID string      `json:"playlist_uploader_id"`
	ChannelURL         string      `json:"channel_url"`
	Resolution         interface{} `json:"resolution"`
	Vcodec             string      `json:"vcodec"`
	NEntries           int         `json:"n_entries"`
	AgeLimit           int         `json:"age_limit"`
}

type DownloadTarget struct {
	URL                		   string
	LatestDownloaded           string
	DownloadMode               string
	Name                       string
	PreferredExtensionForAudio string
	PreferredExtensionForVideo string
	DownloadHistory            []string
	LastChecked                string
	CheckingInterval           string
	Type 					   string
}

type DeleteChannelPayload struct {
	URL string
}

type DeletePlaylistPayload struct {
	URL string
}

type CheckingIntervalPayload struct {
	CheckingInterval string
}

type AddPlaylistPayload struct {
	URL           string
	DownloadMode          string
	FileExtension         string
	DownloadQuality       string
	DownloadEntire bool
	CheckingInterval      string
}

type AddChannelPayload struct {
	URL            string
	DownloadMode          string
	FileExtension         string
	DownloadQuality       string
	DownloadEntire bool
	CheckingInterval      string
}

type DownloadVideoPayload struct {
	 VideoURL 			   string
	 DownloadMode          string
	 FileExtension         string
	 DownloadQuality       string
}

type Response struct {
	Type    string
	Key     string
	Message string
}

type Video struct {
	VideoID string
	DownloadMode string
	DownloadQuality string
	FileExtension string
}
