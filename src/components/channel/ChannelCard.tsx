import Link from "next/link";
import ChannelAvatar from "./ChannelAvatar";

interface ChannelProps {
  avatarUrl: string;
  channelName: string;
  channelHandle: string;
  totalVideos: number;
  totalSize: number;
  checkButton?: boolean;
}

export default function ChannelCard({
  avatarUrl,
  channelName,
  channelHandle,
  totalVideos,
  totalSize,
  checkButton,
}: ChannelProps) {
  const channelUrl = `channels/${channelHandle}`;

  return (
    <div className="flex gap-3 text-white">
      <Link href={channelUrl}>
        <ChannelAvatar avatarUrl={avatarUrl} size={85} />
      </Link>

      <div className="flex flex-col justify-between text-lg">
        <div className="flex flex-col">
          <Link href={channelUrl}>
            <p>{channelName}</p>
          </Link>
          <div className="flex gap-1 text-[#676D75] text-sm">
            <p>{totalVideos} videos</p>
            <p>•</p>
            <p>{formatFileSize(totalSize)}</p>
          </div>
        </div>
        {checkButton && (
          <button className="rounded-full bg-white text-black text-sm py-1 font-semibold w-20">
            Check
          </button>
        )}
      </div>
    </div>
  );
}

function formatFileSize(bytes: number) {
  if (bytes === 0) return "0 Bytes";

  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}
