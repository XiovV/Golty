import VideoMeta from "./VideoMeta";

interface VideoInformationProps {
  videoTitle: string;
  channelName?: string;
  videoSize: string;
  dateDownloaded: string;
}

export default function VideoInformation({
  videoTitle,
  channelName,
  videoSize,
  dateDownloaded,
}: VideoInformationProps) {
  return (
    <div className="flex flex-col gap-1 text-xs max-w-[80%]">
      <p className="text-white text-xs">{videoTitle}</p>

      <VideoMeta
        channelName={channelName ? channelName : ""}
        videoSize={videoSize}
        dateDownloaded={dateDownloaded}
      />
    </div>
  );
}
