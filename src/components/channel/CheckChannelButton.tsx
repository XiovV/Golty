"use client";
import { API_URL } from "@/app/const";
import { useCheckChannel } from "@/hooks/channel/useCheckChannel";
import { IconContext } from "react-icons";
import { LuRefreshCw } from "react-icons/lu";
import { clsx } from "clsx";

interface CheckChannelButtonProps {
  channelId: string;
}

export default function CheckChannelButton({
  channelId,
}: CheckChannelButtonProps) {
  const { loading, checkChannel } = useCheckChannel();

  return (
    <div
      className="hover:cursor-pointer"
      onClick={() => checkChannel(channelId)}
    >
      <IconContext.Provider value={{ size: "1.5em" }}>
        <LuRefreshCw className={clsx({ "animate-spin": loading })} />
      </IconContext.Provider>
    </div>
  );
}
