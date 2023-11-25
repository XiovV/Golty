import Channel from "./channel";
import { Channel as IChannel } from "../types/channel";

interface ChannelListProps {
  channels: IChannel[];
}

export default function ChannelList({ channels }: ChannelListProps) {
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
          />
        );
      })}
    </div>
  );
}
