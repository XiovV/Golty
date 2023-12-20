"use client";
import { fetchChannels } from "@/services/api/channels";
import ChannelCard from "./ChannelCard";
import { API_URL } from "@/app/const";

export default function ChannelList() {
  const { channels, loading } = fetchChannels();

  if (loading && !channels) {
    return <div>Loading...</div>;
  }

  if (!channels) {
    return <div></div>;
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
