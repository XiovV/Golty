"use client";
import { useToast } from "../../components/ui/use-toast";
import { API_URL } from "@/app/const";
import {
  ChannelInfo,
  AddChannelRequest,
  ErrorResponse,
} from "@/hooks/channel/types";

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
        channelInput: formData.get("channelInput")!.toString(),
        channelName: channelInfo.uploader,
        channelHandle: channelInfo.uploaderId,
        avatarUrl: channelInfo.avatarUrl,
      },
      downloadSettings: {
        downloadVideo: Boolean(formData.get("video")),
        downloadAudio: Boolean(formData.get("audio")),
        quality: formData.get("quality")!.toString(),
        format: formData.get("format")!.toString(),
        downloadNewUploads: Boolean(formData.get("downloadAutomatically")),
        downloadEntire: Boolean(formData.get("downloadEntireChannel")),
      },
    };

    const res = await fetch(`${API_URL}/channels`, {
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
