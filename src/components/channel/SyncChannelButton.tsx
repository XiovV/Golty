"use client";
import { useSyncChannel } from "@/hooks/channel";
import { IconContext } from "react-icons";
import { LuRefreshCw } from "react-icons/lu";
import { clsx } from "clsx";

interface CheckChannelButtonProps {
  channelId: string;
}

export default function CheckChannelButton({
  channelId,
}: CheckChannelButtonProps) {
  const { loading, syncChannel } = useSyncChannel();

  return (
    <div
      className="hover:cursor-pointer"
      onClick={() => syncChannel(channelId)}
    >
      <IconContext.Provider value={{ size: "1.5em" }}>
        <LuRefreshCw className={clsx({ "animate-spin": loading })} />
      </IconContext.Provider>
    </div>
  );
}
