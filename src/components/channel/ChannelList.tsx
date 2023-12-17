import ChannelCard from "./ChannelCard";
import { Channel } from "../../types/channel";

async function fetchChannels() {
  const res = await fetch("http://localhost:8080/v1/channels");

  const channels: Channel[] = await res.json();

  return channels;
}

export default async function ChannelList() {
  const channels = await fetchChannels();

  return (
    <div className="flex flex-col gap-6 lg:flex-row lg:flex-wrap lg:gap-x-12">
      {channels.map((channel) => {
        return (
          <>
            <ChannelCard
              key={channel.channelName}
              avatarUrl={channel.avatarUrl}
              channelName={channel.channelName}
              channelHandle={channel.channelHandle}
              totalVideos={channel.totalVideos}
              totalSize={channel.totalSize}
              checkButton
            />
          </>
        );
      })}
    </div>
  );
}
