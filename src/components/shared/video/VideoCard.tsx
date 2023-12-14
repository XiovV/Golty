import Image from "next/image";
import ChannelAvatar from "@/components/channel/ChannelAvatar";
import { formatFileSize } from "@/utils/format";

interface VideoCardProps {
  thumbnailUrl: string;
  title: string;
  avatar?: string;
  channelName?: string;
  videoSize: number;
  dateDownloaded: number;
  showAvatar?: boolean;
}

export default function VideoCard({
  thumbnailUrl,
  title,
  avatar,
  channelName,
  videoSize,
  dateDownloaded,
}: VideoCardProps) {
  return (
    <div className="flex flex-col gap-3 w-[350px]">
      <VideoThumbnail thumbnailUrl={thumbnailUrl} />

      <div className="flex gap-3 items-start">
        {avatar && <ChannelAvatar avatarUrl={avatar} size={40} />}

        <VideoInformation
          videoTitle={title}
          channelName={channelName}
          videoSize={videoSize}
          dateDownloaded={dateDownloaded}
        />
      </div>
    </div>
  );
}

interface VideoInformationProps {
  videoTitle: string;
  channelName?: string;
  videoSize: number;
  dateDownloaded: number;
}

function VideoInformation({
  videoTitle,
  channelName = "",
  videoSize,
  dateDownloaded,
}: VideoInformationProps) {
  return (
    <div className="flex flex-col gap-1 text-xs max-w-[80%]">
      <p className="text-white text-sm">{videoTitle}</p>

      <VideoMeta
        channelName={channelName}
        videoSize={videoSize}
        dateDownloaded={dateDownloaded}
      />
    </div>
  );
}

interface VideoMetaProps {
  channelName?: string;
  videoSize: number;
  dateDownloaded: number;
}

function VideoMeta({ channelName, videoSize, dateDownloaded }: VideoMetaProps) {
  return (
    <div className="flex gap-1 text-[#676D75] text-sm">
      {channelName && (
        <>
          <p>@{channelName}</p>
          <p>•</p>
        </>
      )}
      <p>{formatFileSize(videoSize)}</p>
      <p>•</p>
      <p>{dateDownloaded}</p>
    </div>
  );
}

interface VideoThumbnailProps {
  thumbnailUrl: string;
  width?: number;
  height?: number;
}

function VideoThumbnail({ thumbnailUrl, width, height }: VideoThumbnailProps) {
  return (
    <Image
      priority
      src={thumbnailUrl}
      height={height ? height : 0}
      width={width ? width : 350}
      alt={"video thumbnail"}
    />
  );
}
