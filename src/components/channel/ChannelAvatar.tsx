import Image from "next/image";

interface ChannelAvatarProps {
  avatarUrl: string;
  size: number;
}

export default function ChannelAvatar({ avatarUrl, size }: ChannelAvatarProps) {
  return (
    <div className="rounded-full">
      <Image
        className="rounded-full"
        priority
        src={avatarUrl}
        width={size}
        height={size}
        alt="channel avatar"
      />
    </div>
  );
}
