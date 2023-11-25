import Channel from "./channel";
import { Channel as IChannel } from "../types/channel";

async function fetchChannels(): Promise<IChannel[]> {
  const res = await fetch(`${process.env.API_URL}/channels`, {
    cache: "no-store",
  });

  return await res.json();
}

export default async function ChannelList() {
  const channels = await fetchChannels();
  return (
    <div className="flex flex-col gap-6 lg:flex-row lg:flex-wrap lg:gap-x-12">
      {channels.map((channel) => {
        return (
          <Channel
            key={channel.name}
            avatarUrl={channel.avatarUrl}
            name={channel.name}
            totalVideos={channel.totalVideos}
            totalSize={channel.totalSize}
            checkButton
          />
        );
      })}
    </div>
  );
}
