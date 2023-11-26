import { Video } from "@/types/video";
import VideoCard from "./VideoCard";

interface VideosListProps {
  channelVideosResponse: Promise<Video[]>;
}

export default async function VideosList({
  channelVideosResponse,
}: VideosListProps) {
  const videos = await channelVideosResponse;
  return (
    <div className="flex flex-col gap-3 flex-wrap">
      {videos.map((video) => {
        return (
          <VideoCard
            thumbnailUrl={video.thumbnailUrl}
            title={video.title}
            videoSize={video.size}
            dateDownloaded={video.dateDownloaded}
          />
        );
      })}
    </div>
  );
}
