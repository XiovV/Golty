package main

import "os"

var (
	API_KEY               = os.Getenv("API_KEY")
	API_ENDPOINT_ID       = "https://www.googleapis.com/youtube/v3/search?part=snippet&channelId="
	API_ENDPOINT_NAME     = "https://www.googleapis.com/youtube/v3/channels?part=contentDetails&forUsername="
	API_ENDPOINT_PLAYLIST = "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId="
	MAX_RESULTS           = "maxResults=1"
	ORDER_BY              = "order=date"
	TYPE                  = "type=video"
)
