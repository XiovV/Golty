"use client";
import { fetchChannels } from "@/services/api/channels";
import ChannelCard from "./ChannelCard";
import { useEffect } from "react";

const POLLING_RATE = 3000;

export default function ChannelList() {
  const { channels, loading } = fetchChannels();

  // useEffect(() => {
  //   const intervalId = setInterval(() => {
  //     fetchData();
  //   }, POLLING_RATE);
  //
  //   return () => clearInterval(intervalId);
  // }, []);

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
            avatarUrl={channel.avatarUrl}
            channelName={channel.channelName}
            channelHandle={channel.channelHandle}
            totalVideos={channel.totalVideos}
            totalSize={channel.totalSize}
            downloading={channel.state === "downloading"}
            checkButton
          />
        );
      })}
    </div>
  );
}
