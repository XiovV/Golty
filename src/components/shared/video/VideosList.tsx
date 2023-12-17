import VideoCard from "./VideoCard";

interface VideosListProps {
  channelHandle: string;
}

interface Video {
  videoId: string;
  title: string;
  thumbnailUrl: string;
  size: number;
  downloadDate: number;
  duration: string;
}

async function fetchVideos(channelHandle: string) {
  const res = await fetch(
    `http://localhost:8080/v1/channels/videos/${channelHandle}`,
    { cache: "no-cache" },
  );

  const videos: Video[] = await res.json();

  return videos;
}

export default async function VideosList({ channelHandle }: VideosListProps) {
  const videos = await fetchVideos(channelHandle);

  return (
    <div className="flex flex-col gap-5 lg:flex-row lg:flex-wrap content-start items-start">
      {videos.map((video) => {
        return (
          <VideoCard
            key={video.videoId}
            thumbnailUrl={video.thumbnailUrl}
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
