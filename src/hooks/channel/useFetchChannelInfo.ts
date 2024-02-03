import { useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import { ChannelInfo } from "@/hooks/channel/types";
import { API_URL } from "@/app/const/index";

export const useFetchChannelInfo = () => {
  const [loading, setLoading] = useState(false);
  const [channelInfo, setChannelInfo] = useState<ChannelInfo>();

  const getChannelInfo = useDebouncedCallback(async (channelInput: string) => {
    if (!channelInput) {
      return;
    }

    setChannelInfo(undefined);
    setLoading(true);
    const res = await fetch(`${API_URL}/channels/info/${channelInput}`, {
      cache: "no-cache",
    });

    if (res.status !== 200) {
      setLoading(false);
      return;
    }

    const channelInfo: ChannelInfo = await res.json();
    setLoading(false);
    setChannelInfo(channelInfo);
  }, 500);

  return { loading, channelInfo, getChannelInfo };
};
