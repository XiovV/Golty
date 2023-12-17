import Image from "next/image";
import ChannelAvatar from "@/components/channel/ChannelAvatar";
import { formatFileSize, formatTimeAgo } from "@/utils/format";

interface VideoCardProps {
  thumbnailUrl: string;
  title: string;
  avatar?: string;
  channelName?: string;
  videoSize: number;
  downloadDate: number;
  showAvatar?: boolean;
  duration: string;
}

export default function VideoCard({
  thumbnailUrl,
  title,
  avatar,
  channelName,
  videoSize,
  downloadDate,
  duration,
}: VideoCardProps) {
  return (
    <div className="flex flex-col gap-3 w-[350px]">
      <VideoThumbnail thumbnailUrl={thumbnailUrl} duration={duration} />

      <div className="flex gap-3 items-start">
        {avatar && <ChannelAvatar avatarUrl={avatar} size={40} />}

        <VideoInformation
          videoTitle={title}
          channelName={channelName}
          videoSize={videoSize}
          downloadDate={downloadDate}
        />
      </div>
    </div>
  );
}

interface VideoInformationProps {
  videoTitle: string;
  channelName?: string;
  videoSize: number;
  downloadDate: number;
}

function VideoInformation({
  videoTitle,
  channelName = "",
  videoSize,
  downloadDate,
}: VideoInformationProps) {
  return (
    <div className="flex flex-col gap-1 text-xs max-w-[80%]">
      <p className="text-white text-sm">{videoTitle}</p>

      <VideoMeta
        channelName={channelName}
        videoSize={videoSize}
        downloadDate={downloadDate}
      />
    </div>
  );
}

interface VideoMetaProps {
  channelName?: string;
  videoSize: number;
  downloadDate: number;
}

function VideoMeta({ channelName, videoSize, downloadDate }: VideoMetaProps) {
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
      <p>{formatTimeAgo(downloadDate)}</p>
    </div>
  );
}

interface VideoThumbnailProps {
  thumbnailUrl: string;
  duration: string;
  width?: number;
  height?: number;
}

function VideoThumbnail({
  thumbnailUrl,
  width,
  height,
  duration,
}: VideoThumbnailProps) {
  return (
    <div className="relative">
      <Image
        priority
        src={thumbnailUrl}
        height={height ? height : 0}
        width={width ? width : 350}
        alt={"video thumbnail"}
      />

      <div className="absolute bottom-0 right-0 p-2 text-white">
        <span className="rounded-md bg-black p-1 text-xs ">{duration}</span>
      </div>
    </div>
  );
}
