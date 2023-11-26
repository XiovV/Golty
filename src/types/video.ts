export interface ChannelVideosResponse {
  videos: Video[];
}

export interface Video {
  title: string;
  size: string;
  thumbnailUrl: string;
  channelName: string;
  dateDownloaded: string;
}