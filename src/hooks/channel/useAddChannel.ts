"use client";
import { useToast } from "../../components/ui/use-toast";
import { useRouter } from "next/navigation";
import { API_URL } from "@/app/const";
import {
  ChannelInfo,
  AddChannelRequest,
  ErrorResponse,
} from "@/hooks/channel/types";

export const useAddChannel = () => {
  const { toast } = useToast();
  const router = useRouter();

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
        avatarUrl: channelInfo.avatar_url,
      },
      downloadSettings: {
        downloadVideo: Boolean(formData.get("video")),
        downloadAudio: Boolean(formData.get("audio")),
        resolution: formData.get("resolution")!.toString(),
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
