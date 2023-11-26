import ChannelCard from "./ChannelCard";
import { Channel as IChannel } from "../../types/channel";

interface ChannelListProps {
  channelsResponse: Promise<IChannel[]>;
}

export default async function ChannelList({
  channelsResponse,
}: ChannelListProps) {
  const channels = await channelsResponse;

  return (
    <div className="flex flex-col gap-6 lg:flex-row lg:flex-wrap lg:gap-x-12">
      {channels.map((channel) => {
        return (
          <ChannelCard
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
