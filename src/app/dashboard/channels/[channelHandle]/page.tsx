import { Channel } from "@/services/api/channels";
import { API_URL } from "@/app/const";
import CheckChannelButton from "@/components/channel/CheckChannelButton";
import {
  TopBar,
  DesktopButtons,
  MobileButtons,
} from "@/components/navigation/TopBar";

import ChannelCard from "@/components/channel/ChannelCard";
import VideosList from "@/components/shared/video/VideosList";

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
      <TopBar title={channel.channelName}>
        <DesktopButtons>
          <CheckChannelButton channelId={channel.id} />
        </DesktopButtons>

        <MobileButtons title={channel.channelName}>
          <CheckChannelButton channelId={channel.id} />
        </MobileButtons>
      </TopBar>

      <div className="flex flex-col gap-8 m-5">
        <ChannelCard
          avatarUrl={`${API_URL}/assets/${channel.avatarUrl}`}
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
