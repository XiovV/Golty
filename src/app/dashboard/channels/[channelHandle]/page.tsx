import { API_URL } from "@/app/const";
import SyncChannelButton from "@/components/channel/SyncChannelButton";
import DeleteChannelButton from "@/components/channel/DeleteChannelButton";
import {
  TopBar,
  DesktopButtons,
  MobileButtons,
} from "@/components/navigation/TopBar";

import ChannelCard from "@/components/channel/ChannelCard";
import VideosList from "@/components/shared/video/VideosList";
import { Channel } from "@/hooks/channel/types";

async function fetchChannel(channelHandle: string): Promise<Channel> {
  const res = await fetch(`${API_URL}/channels/${channelHandle}`, {
    cache: "no-store",
  });

  const channel: Channel = await res.json();

  return channel;
}

export default async function Page({
  params,
}: {
  params: { channelHandle: string };
}) {
  const channel = await fetchChannel(params.channelHandle);

  return (
    <main>
      <TopBar>
        <DesktopButtons>
          <SyncChannelButton channelId={channel.id} />
          <DeleteChannelButton channelId={channel.id} />
        </DesktopButtons>

        <MobileButtons title={channel.channelName}>
          <SyncChannelButton channelId={channel.id} />
          <DeleteChannelButton channelId={channel.id} />
        </MobileButtons>
      </TopBar>

      <div className="flex flex-col gap-8 m-5">
        <ChannelCard
          avatarUrl={`${API_URL}/assets/${channel.avatarUrl}`}
          channelId={channel.id}
          channelName={channel.channelName}
          channelHandle={channel.channelHandle}
          totalVideos={channel.totalVideos}
          totalSize={channel.totalSize}
        />

        <div className="mx-auto lg:mx-5 pb-16">
          <VideosList channelId={channel.id} />
        </div>
      </div>
    </main>
  );
}
