export interface Channel {
  id: string;
  channelName: string;
  channelHandle: string;
  channelUrl: string;
  avatarUrl: string;
  totalVideos: number;
  totalSize: number;
  state: string;
}

export interface ChannelInfo {
  uploader_url: string;
  uploader_id: string;
  uploader: string;
  avatar_url: string;
}

export interface ErrorResponse {
  message: string;
}

export interface AddChannelRequest {
  channel: {
    channelUrl: string;
    channelName: string;
    channelHandle: string;
    avatarUrl: string;
  };
  downloadSettings: {
    downloadVideo: boolean;
    downloadAudio: boolean;
    format: string;
    resolution: string;
    downloadNewUploads: boolean;
    downloadEntire: boolean;
  };
}

export interface SyncChannelResponse {
  missingVideos: number;
}
