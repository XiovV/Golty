"use client";
import Link from "next/link";
import ChannelAvatar from "./ChannelAvatar";
import { formatFileSize } from "@/utils/format";
import PulseLoader from "react-spinners/PulseLoader";

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
  checkButton,
  downloading = false,
}: ChannelProps) {
  const channelUrl = `channels/${channelHandle}`;

  return (
    <div className="flex gap-3 text-white">
      <Link href={channelUrl}>
        <ChannelAvatar avatarUrl={avatarUrl} size={85} />
      </Link>

      <div className="flex flex-col justify-between text-lg">
        <div className="flex flex-col">
          <div className="flex gap-2 items-center">
            <Link href={channelUrl}>
              <p>{channelName}</p>
            </Link>

            <PulseLoader
              color={"#ffffff"}
              loading={downloading}
              size={7}
              aria-label="Loading Spinner"
              data-testid="loader"
            />
          </div>

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
