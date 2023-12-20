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

    let toastDescription = "No new videos detected.";

    if (syncChannelResponse.missingVideos == 1) {
      toastDescription = "1 new video detected.";
    } else if (syncChannelResponse.missingVideos > 1) {
      toastDescription = `${syncChannelResponse.missingVideos} new videos detected.`;
    }

    toast({
      title: "Checking for new videos completed.",
      description: toastDescription,
    });
  };

  return { loading, syncChannel };
};
