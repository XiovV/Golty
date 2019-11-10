package main

import "time"

type ChannelBody struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string `json:"channelTitle"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
}

type NameBody struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind    string `json:"kind"`
		Etag    string `json:"etag"`
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
				Standard struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"standard"`
				Maxres struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"maxres"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
			PlaylistID   string `json:"playlistId"`
			Position     int    `json:"position"`
			ResourceID   struct {
				Kind    string `json:"kind"`
				VideoID string `json:"videoId"`
			} `json:"resourceId"`
		} `json:"snippet"`
	} `json:"items"`
}

type UserBody struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind           string `json:"kind"`
		Etag           string `json:"etag"`
		ID             string `json:"id"`
		ContentDetails struct {
			RelatedPlaylists struct {
				Uploads      string `json:"uploads"`
				WatchHistory string `json:"watchHistory"`
				WatchLater   string `json:"watchLater"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type UploadID struct {
	ChannelURL string
	UploadsID  string
}

type FailedVideos struct {
	VideoID string
}

type Channel struct {
	ChannelURL       string
	LatestDownloaded string
}

type Payload struct {
	ChannelURL string
}

type Response struct {
	Type    string
	Key     string
	Message string
}

type Option struct {
	Resume bool // resume failed or cancelled download
	Rename bool // rename output file using video title
	Mp3    bool // extract audio using ffmpeg
}

type Video struct {
	Id, Title, Author, Keywords, Thumbnail_url string
	Avg_rating                                 float32
	View_count, Length_seconds                 int
	Formats                                    []Format
	Filename                                   string
}

type Format struct {
	Itag                     int
	Video_type, Quality, Url string
}

type PlayerResponse struct {
	ResponseContext struct {
		ServiceTrackingParams []struct {
			Service string `json:"service"`
			Params  []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"params"`
		} `json:"serviceTrackingParams"`
	} `json:"responseContext"`
	PlayabilityStatus struct {
		Status          string `json:"status"`
		PlayableInEmbed bool   `json:"playableInEmbed"`
	} `json:"playabilityStatus"`
	StreamingData struct {
		ExpiresInSeconds string `json:"expiresInSeconds"`
		Formats          []struct {
			Itag             int    `json:"itag"`
			URL              string `json:"url"`
			MimeType         string `json:"mimeType"`
			Bitrate          int    `json:"bitrate"`
			Width            int    `json:"width"`
			Height           int    `json:"height"`
			LastModified     string `json:"lastModified"`
			ContentLength    string `json:"contentLength"`
			Quality          string `json:"quality"`
			QualityLabel     string `json:"qualityLabel"`
			ProjectionType   string `json:"projectionType"`
			AverageBitrate   int    `json:"averageBitrate,omitempty"`
			AudioQuality     string `json:"audioQuality"`
			ApproxDurationMs string `json:"approxDurationMs,omitempty"`
			AudioSampleRate  string `json:"audioSampleRate,omitempty"`
			AudioChannels    int    `json:"audioChannels,omitempty"`
		} `json:"formats"`
		AdaptiveFormats []struct {
			Itag      int    `json:"itag"`
			URL       string `json:"url"`
			MimeType  string `json:"mimeType"`
			Bitrate   int    `json:"bitrate"`
			Width     int    `json:"width,omitempty"`
			Height    int    `json:"height,omitempty"`
			InitRange struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"initRange"`
			IndexRange struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"indexRange"`
			LastModified     string `json:"lastModified"`
			ContentLength    string `json:"contentLength"`
			Quality          string `json:"quality"`
			Fps              int    `json:"fps,omitempty"`
			QualityLabel     string `json:"qualityLabel,omitempty"`
			ProjectionType   string `json:"projectionType"`
			AverageBitrate   int    `json:"averageBitrate"`
			ApproxDurationMs string `json:"approxDurationMs"`
			HighReplication  bool   `json:"highReplication,omitempty"`
			AudioQuality     string `json:"audioQuality,omitempty"`
			AudioSampleRate  string `json:"audioSampleRate,omitempty"`
			AudioChannels    int    `json:"audioChannels,omitempty"`
		} `json:"adaptiveFormats"`
	} `json:"streamingData"`
	PlaybackTracking struct {
		VideostatsPlaybackURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsPlaybackUrl"`
		VideostatsDelayplayURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsDelayplayUrl"`
		VideostatsWatchtimeURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videostatsWatchtimeUrl"`
		PtrackingURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"ptrackingUrl"`
		QoeURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"qoeUrl"`
		SetAwesomeURL struct {
			BaseURL                 string `json:"baseUrl"`
			ElapsedMediaTimeSeconds int    `json:"elapsedMediaTimeSeconds"`
		} `json:"setAwesomeUrl"`
		AtrURL struct {
			BaseURL                 string `json:"baseUrl"`
			ElapsedMediaTimeSeconds int    `json:"elapsedMediaTimeSeconds"`
		} `json:"atrUrl"`
	} `json:"playbackTracking"`
	Captions struct {
		PlayerCaptionsRenderer struct {
			BaseURL    string `json:"baseUrl"`
			Visibility string `json:"visibility"`
		} `json:"playerCaptionsRenderer"`
		PlayerCaptionsTracklistRenderer struct {
			CaptionTracks []struct {
				BaseURL string `json:"baseUrl"`
				Name    struct {
					SimpleText string `json:"simpleText"`
				} `json:"name"`
				VssID          string `json:"vssId"`
				LanguageCode   string `json:"languageCode"`
				Kind           string `json:"kind"`
				IsTranslatable bool   `json:"isTranslatable"`
			} `json:"captionTracks"`
			AudioTracks []struct {
				CaptionTrackIndices []int  `json:"captionTrackIndices"`
				Visibility          string `json:"visibility"`
			} `json:"audioTracks"`
			TranslationLanguages []struct {
				LanguageCode string `json:"languageCode"`
				LanguageName struct {
					SimpleText string `json:"simpleText"`
				} `json:"languageName"`
			} `json:"translationLanguages"`
			DefaultAudioTrackIndex int `json:"defaultAudioTrackIndex"`
		} `json:"playerCaptionsTracklistRenderer"`
	} `json:"captions"`
	VideoDetails struct {
		VideoID          string   `json:"videoId"`
		Title            string   `json:"title"`
		LengthSeconds    string   `json:"lengthSeconds"`
		Keywords         []string `json:"keywords"`
		ChannelID        string   `json:"channelId"`
		IsOwnerViewing   bool     `json:"isOwnerViewing"`
		ShortDescription string   `json:"shortDescription"`
		IsCrawlable      bool     `json:"isCrawlable"`
		Thumbnail        struct {
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
		} `json:"thumbnail"`
		UseCipher         bool    `json:"useCipher"`
		AverageRating     float64 `json:"averageRating"`
		AllowRatings      bool    `json:"allowRatings"`
		ViewCount         string  `json:"viewCount"`
		Author            string  `json:"author"`
		IsPrivate         bool    `json:"isPrivate"`
		IsUnpluggedCorpus bool    `json:"isUnpluggedCorpus"`
		IsLiveContent     bool    `json:"isLiveContent"`
	} `json:"videoDetails"`
	Annotations []struct {
		PlayerAnnotationsUrlsRenderer struct {
			InvideoURL         string `json:"invideoUrl"`
			LoadPolicy         string `json:"loadPolicy"`
			AllowInPlaceSwitch bool   `json:"allowInPlaceSwitch"`
		} `json:"playerAnnotationsUrlsRenderer"`
	} `json:"annotations"`
	PlayerConfig struct {
		AudioConfig struct {
			LoudnessDb           float64 `json:"loudnessDb"`
			PerceptualLoudnessDb float64 `json:"perceptualLoudnessDb"`
		} `json:"audioConfig"`
		StreamSelectionConfig struct {
			MaxBitrate string `json:"maxBitrate"`
		} `json:"streamSelectionConfig"`
		MediaCommonConfig struct {
			DynamicReadaheadConfig struct {
				MaxReadAheadMediaTimeMs int `json:"maxReadAheadMediaTimeMs"`
				MinReadAheadMediaTimeMs int `json:"minReadAheadMediaTimeMs"`
				ReadAheadGrowthRateMs   int `json:"readAheadGrowthRateMs"`
			} `json:"dynamicReadaheadConfig"`
		} `json:"mediaCommonConfig"`
	} `json:"playerConfig"`
	Storyboards struct {
		PlayerStoryboardSpecRenderer struct {
			Spec string `json:"spec"`
		} `json:"playerStoryboardSpecRenderer"`
	} `json:"storyboards"`
	Microformat struct {
		PlayerMicroformatRenderer struct {
			Thumbnail struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			Embed struct {
				IframeURL      string `json:"iframeUrl"`
				FlashURL       string `json:"flashUrl"`
				Width          int    `json:"width"`
				Height         int    `json:"height"`
				FlashSecureURL string `json:"flashSecureUrl"`
			} `json:"embed"`
			Title struct {
				SimpleText string `json:"simpleText"`
			} `json:"title"`
			Description struct {
				SimpleText string `json:"simpleText"`
			} `json:"description"`
			LengthSeconds        string   `json:"lengthSeconds"`
			OwnerProfileURL      string   `json:"ownerProfileUrl"`
			OwnerGplusProfileURL string   `json:"ownerGplusProfileUrl"`
			ExternalChannelID    string   `json:"externalChannelId"`
			AvailableCountries   []string `json:"availableCountries"`
			IsUnlisted           bool     `json:"isUnlisted"`
			HasYpcMetadata       bool     `json:"hasYpcMetadata"`
			ViewCount            string   `json:"viewCount"`
			Category             string   `json:"category"`
			PublishDate          string   `json:"publishDate"`
			OwnerChannelName     string   `json:"ownerChannelName"`
			UploadDate           string   `json:"uploadDate"`
		} `json:"playerMicroformatRenderer"`
	} `json:"microformat"`
	TrackingParams string `json:"trackingParams"`
	Attestation    struct {
		PlayerAttestationRenderer struct {
			Challenge    string `json:"challenge"`
			BotguardData struct {
				Program        string `json:"program"`
				InterpreterURL string `json:"interpreterUrl"`
			} `json:"botguardData"`
		} `json:"playerAttestationRenderer"`
	} `json:"attestation"`
	Messages []struct {
		MealbarPromoRenderer struct {
			MessageTexts []struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"messageTexts"`
			ActionButton struct {
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						URLEndpoint         struct {
							URL    string `json:"url"`
							Target string `json:"target"`
						} `json:"urlEndpoint"`
					} `json:"navigationEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"buttonRenderer"`
			} `json:"actionButton"`
			DismissButton struct {
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					ServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						FeedbackEndpoint    struct {
							FeedbackToken string `json:"feedbackToken"`
							UIActions     struct {
								HideEnclosingContainer bool `json:"hideEnclosingContainer"`
							} `json:"uiActions"`
						} `json:"feedbackEndpoint"`
					} `json:"serviceEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"buttonRenderer"`
			} `json:"dismissButton"`
			TriggerCondition    string `json:"triggerCondition"`
			Style               string `json:"style"`
			TrackingParams      string `json:"trackingParams"`
			ImpressionEndpoints []struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				FeedbackEndpoint    struct {
					FeedbackToken string `json:"feedbackToken"`
					UIActions     struct {
						HideEnclosingContainer bool `json:"hideEnclosingContainer"`
					} `json:"uiActions"`
				} `json:"feedbackEndpoint"`
			} `json:"impressionEndpoints"`
			IsVisible    bool `json:"isVisible"`
			MessageTitle struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"messageTitle"`
		} `json:"mealbarPromoRenderer"`
	} `json:"messages"`
	AdSafetyReason struct {
		ApmUserPreference struct {
		} `json:"apmUserPreference"`
		IsEmbed bool `json:"isEmbed"`
	} `json:"adSafetyReason"`
}
