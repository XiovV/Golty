import { API_URL } from "@/app/const";
import VideoCard from "./VideoCard";

interface VideosListProps {
  channelId: string;
}

interface Video {
  videoId: string;
  title: string;
  thumbnailUrl: string;
  size: number;
  downloadDate: number;
  duration: string;
}

async function fetchVideos(channelId: string) {
  const res = await fetch(`${API_URL}/channels/videos/${channelId}`, {
    cache: "no-cache",
  });

  const videos: Video[] = await res.json();

  return videos;
}

export default async function VideosList({ channelId }: VideosListProps) {
  const videos = await fetchVideos(channelId);

  return (
    <div className="flex flex-col gap-5 lg:flex-row lg:flex-wrap content-start items-start">
      {videos.map((video) => {
        return (
          <VideoCard
            key={video.videoId}
            thumbnailUrl={`${API_URL}/assets/${video.thumbnailUrl}`}
            title={video.title}
            videoSize={video.size}
            downloadDate={video.downloadDate}
            duration={video.duration}
          />
        );
      })}
    </div>
  );
}
