import Image from "next/image";
import Link from "next/link";

interface ChannelProps {
  avatarUrl: string;
  name: string;
  totalVideos: number;
  totalSize: string;
  checkButton?: boolean;
}

export default function Channel({
  avatarUrl,
  name,
  totalVideos,
  totalSize,
  checkButton,
}: ChannelProps) {
  const channelUrl = `channels/${name}`;

  return (
    <div className="flex gap-3 text-white">
      <Link href={channelUrl}>
        <Image priority alt="" src={avatarUrl} width={85} height={85} />
      </Link>

      <div className="flex flex-col justify-between text-lg">
        <div className="flex flex-col">
          <Link href={channelUrl}>
            <p>{name}</p>
          </Link>
          <div className="flex gap-1 text-[#676D75] text-sm">
            <p>{totalVideos} videos</p>
            <p>•</p>
            <p>{totalSize}</p>
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