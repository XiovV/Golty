import ChannelAvatar from "./ChannelAvatar";
import Link from "next/link";

interface ChannelInfoCardProps {
  avatarUrl: string;
  channelUrl: string;
  name: string;
  handle: string;
}

export default function ChannelInfoCard({
  avatarUrl,
  channelUrl,
  name,
  handle,
}: ChannelInfoCardProps) {
  return (
    <div className="flex gap-3 text-white">
      <Link href={channelUrl}>
        <ChannelAvatar avatarUrl={avatarUrl} size={85} />
      </Link>

      <div className="flex flex-col justify-between text-lg">
        <div className="flex flex-col gap-1">
          <Link href={channelUrl}>
            <p>{name}</p>
          </Link>
          <Link href={channelUrl}>
            <p className="text-[#676D75] text-sm">{handle}</p>
          </Link>
        </div>
      </div>
    </div>
  );
}
