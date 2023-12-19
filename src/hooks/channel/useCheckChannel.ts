import { API_URL } from "@/app/const/index";
import { useState } from "react";
import { CheckChannelResponse } from "./types";
import { useToast } from "../../components/ui/use-toast";

export const useCheckChannel = () => {
  const [loading, setLoading] = useState(false);
  const { toast } = useToast();

  const checkChannel = async (channelId: string) => {
    setLoading(true);
    const res = await fetch(`${API_URL}/channels/check/${channelId}`);

    const checkChannelResponse: CheckChannelResponse = await res.json();

    setLoading(false);

    toast({
      title: "Checking for new videos completed.",
      description:
        checkChannelResponse.missingVideos == 0
          ? "No new videos detected"
          : `${checkChannelResponse.missingVideos} videos detected.`,
    });
  };

  return { loading, checkChannel };
};
