"use client";
import Link from "next/link";
import ChannelAvatar from "./ChannelAvatar";
import { formatFileSize } from "@/utils/format";
import { clsx } from "clsx";

import { LuRefreshCw } from "react-icons/lu";

interface ChannelProps {
  avatarUrl: string;
  channelName: string;
  channelHandle: string;
  totalVideos: number;
  totalSize: number;
  checkButton?: boolean;
  downloading?: boolean;
}

export default function ChannelCard({
  avatarUrl,
  channelName,
  channelHandle,
  totalVideos,
  totalSize,
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
            <LuRefreshCw className={clsx({ "animate-spin": false })} />
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
