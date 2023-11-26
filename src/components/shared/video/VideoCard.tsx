import Image from "next/image";
import VideoThumbnail from "./VideoThumbnail";
import VideoInformation from "./VideoInformation";
import ChannelAvatar from "@/components/channel/ChannelAvatar";

interface VideoCardProps {
  thumbnailUrl: string;
  title: string;
  avatar?: string;
  channelName?: string;
  videoSize: string;
  dateDownloaded: string;
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
