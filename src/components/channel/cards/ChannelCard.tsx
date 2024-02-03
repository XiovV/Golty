"use client";
import Link from "next/link";
import ChannelAvatar from "../ChannelAvatar";
import { formatFileSize } from "@/utils/format";
import SyncChannelButton from "../buttons/SyncChannelButton";

interface ChannelProps {
  avatarUrl: string;
  channelId: string;
  channelName: string;
  channelHandle: string;
  totalVideos: number;
  totalSize: number;
  syncButton?: boolean;
}

export default function ChannelCard({
  avatarUrl,
  channelId,
  channelName,
  channelHandle,
  totalVideos,
  totalSize,
  syncButton,
}: ChannelProps) {
  const channelUrl = `channels/${channelHandle}`;

  return (
    <div className="flex gap-3 text-white">
      <Link href={channelUrl}>
        <ChannelAvatar avatarUrl={avatarUrl} size={85} />
      </Link>

      <div className="flex flex-col justify-between text-lg">
        <div className="flex flex-col">
          <div className="flex gap-3 items-center">
            <Link href={channelUrl}>
              <p>{channelName}</p>
            </Link>
            {syncButton && (
              <SyncChannelButton channelId={channelId} size="1em" />
            )}
          </div>

          <p className="text-[#676D75] text-sm mt-1">{channelHandle}</p>

          <div className="flex gap-1 text-[#676D75] text-sm mt-2">
            <p>{totalVideos} videos</p>
            <p>•</p>
            <p>{formatFileSize(totalSize)}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
