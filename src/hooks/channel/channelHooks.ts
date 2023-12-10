import { useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import { useToast } from "../../components/ui/use-toast";

interface ChannelInfo {
  uploader_id: string;
  uploader: string;
  avatar: Avatar;
}

interface Avatar {
  url: string;
}

interface ErrorResponse {
  message: string;
}

interface AddChannelRequest {
  channel: {
    channelUrl: string;
    channelName: string;
    channelHandle: string;
    avatarUrl: string;
  };
  downloadSettings: {
    downloadVideo: boolean;
    downloadAudio: boolean;
    format: string;
    resolution: string;
    downloadAutomatically: boolean;
    downloadEntireChannel: boolean;
  };
}

export const useAddChannel = () => {
  const { toast } = useToast();

  const addChannel = async (
    e: React.FormEvent<HTMLFormElement>,
    channelInfo: ChannelInfo,
  ) => {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);

    const body: AddChannelRequest = {
      channel: {
        channelUrl: formData.get("channelUrl")!.toString(),
        channelName: channelInfo.uploader,
        channelHandle: channelInfo.uploader_id,
        avatarUrl: channelInfo.avatar.url,
      },
      downloadSettings: {
        downloadVideo: Boolean(formData.get("video")),
        downloadAudio: Boolean(formData.get("audio")),
        resolution: formData.get("resolution")!.toString(),
        format: formData.get("format")!.toString(),
        downloadAutomatically: Boolean(formData.get("downloadAutomatically")),
        downloadEntireChannel: Boolean(formData.get("downloadEntireChannel")),
      },
    };

    const res = await fetch("http://localhost:8080/v1/channels", {
      method: "POST",
      body: JSON.stringify(body),
      headers: { "Content-Type": "application/json" },
    });

    if (res.status !== 201) {
      const err: ErrorResponse = await res.json();

      toast({
        title: `Unable to add the channel! (${res.status} ${res.statusText})`,
        description: err.message,
      });

      return;
    }

    toast({
      title: `${channelInfo.uploader} added successfully!`,
    });
  };

  return { addChannel };
};

export const useFetchChannelInfo = () => {
  const [loading, setLoading] = useState(false);
  const [channelInfo, setChannelInfo] = useState<ChannelInfo>();

  const getChannelInfo = useDebouncedCallback(async (channelUrl: string) => {
    if (!channelUrl || !channelUrl.includes("https://www.youtube.com/")) {
      return;
    }

    setChannelInfo(undefined);
    setLoading(true);
    const res = await fetch(
      `http://localhost:8080/v1/channels/info/${channelUrl}`,
      { cache: "no-cache" },
    );

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
