"use client";
import ChannelCard from "./ChannelCard";
import { API_URL } from "@/app/const";
import { useFetchChannels } from "@/hooks/channel/useFetchChannels";
import { useEffect } from "react";

export default function ChannelList() {
  const { channels, loading, fetchData } = useFetchChannels();

  useEffect(() => {
    fetchData();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!loading && !channels) {
    return (
      <div>
        No channels here so far! Press the + icon on the top right to add one.
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-6 lg:flex-row lg:flex-wrap lg:gap-x-12">
      {channels.map((channel) => {
        return (
          <ChannelCard
            key={channel.channelName}
            avatarUrl={`${API_URL}/assets/${channel.avatarUrl}`}
            channelId={channel.id}
            channelName={channel.channelName}
            channelHandle={channel.channelHandle}
            totalVideos={channel.totalVideos}
            totalSize={channel.totalSize}
            syncButton
          />
        );
      })}
    </div>
  );
}
