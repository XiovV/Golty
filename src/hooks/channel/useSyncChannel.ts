import { API_URL } from "@/app/const/index";
import { useState } from "react";
import { SyncChannelResponse } from "./types";
import { useToast } from "../../components/ui/use-toast";

export const useSyncChannel = () => {
  const [loading, setLoading] = useState(false);
  const { toast } = useToast();

  const syncChannel = async (channelId: string) => {
    setLoading(true);
    const res = await fetch(`${API_URL}/channels/sync/${channelId}`);

    const syncChannelResponse: SyncChannelResponse = await res.json();

    setLoading(false);

    toast({
      title: "Checking for new videos completed.",
      description:
        syncChannelResponse.missingVideos == 0
          ? "No new videos detected"
          : `${syncChannelResponse.missingVideos} videos detected.`,
    });
  };

  return { loading, syncChannel };
};
