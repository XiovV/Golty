interface VideoMetaProps {
  channelName?: string;
  videoSize: string;
  dateDownloaded: string;
}

export default function VideoMeta({
  channelName,
  videoSize,
  dateDownloaded,
}: VideoMetaProps) {
  return (
    <div className="flex gap-1 text-[#676D75]">
      {channelName && (
        <>
          <p>@{channelName}</p>
          <p>•</p>
        </>
      )}
      <p>{videoSize}</p>
      <p>•</p>
      <p>{dateDownloaded}</p>
    </div>
  );
}
