import Image from "next/image";

interface ChannelProps {
  avatarUrl: string;
  name: string;
  totalVideos: number;
  totalSize: string;
}

export default function Channel({
  avatarUrl,
  name,
  totalVideos,
  totalSize,
}: ChannelProps) {
  return (
    <div className="flex gap-3 text-white">
      <Image priority alt="" src={avatarUrl} width={85} height={100} />

      <div className="flex flex-col justify-between text-lg">
        <div className="flex flex-col">
          <p>{name}</p>
          <div className="flex gap-1 text-[#676D75] text-sm">
            <p>{totalVideos}</p>
            <p>•</p>
            <p>{totalSize}</p>
          </div>
        </div>
        <button className="rounded-full bg-white text-black text-sm py-1 font-semibold">
          Check
        </button>
      </div>
    </div>
  );
}
