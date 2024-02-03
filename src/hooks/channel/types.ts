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
  uploaderUrl: string;
  uploaderId: string;
  uploader: string;
  avatarUrl: string;
}

export interface ErrorResponse {
  message: string;
}

export interface AddChannelRequest {
  channel: {
    channelInput: string;
    channelName: string;
    channelHandle: string;
    avatarUrl: string;
  };
  downloadSettings: {
    downloadVideo: boolean;
    downloadAudio: boolean;
    format: string;
    quality: string;
    downloadNewUploads: boolean;
    downloadEntire: boolean;
  };
}

export interface SyncChannelResponse {
  missingVideos: number;
}
