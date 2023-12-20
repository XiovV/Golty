"use client";
import { useSyncChannel } from "@/hooks/channel";
import { IconContext } from "react-icons";
import { LuRefreshCw } from "react-icons/lu";
import { clsx } from "clsx";

interface CheckChannelButtonProps {
  channelId: string;
  size?: string;
}

export default function SyncChannelButton({
  channelId,
  size = "1.5em",
}: CheckChannelButtonProps) {
  const { loading, syncChannel } = useSyncChannel();

  return (
    <div
      className="hover:cursor-pointer"
      onClick={() => syncChannel(channelId)}
    >
      <IconContext.Provider value={{ size: size }}>
        <LuRefreshCw className={clsx({ "animate-spin": loading })} />
      </IconContext.Provider>
    </div>
  );
}
